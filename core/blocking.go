package core

import (
	"context"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/olekukonko/tablewriter"
	"github.com/rfyiamcool/go-netflow"
	"github.com/rfyiamcool/go-netflow/constants"
	"github.com/rfyiamcool/go-netflow/rpc"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var nf netflow.Interface

func Start(pname string) {
	var err error
	// Initialize netflow instance with error handling
	nf, err = netflow.New(netflow.WithName(pname))
	if err != nil {
		log.Fatalf("Failed to create netflow instance: %v", err)
		return
	}
	// Start netflow instance
	err = nf.Start()
	if err != nil {
		log.Fatalf("Failed to core netflow instance: %v", err)
		nf.Stop() // Ensure cleanup if core fails
		return
	}

	// Ensure resources are cleaned up on function exit
	defer func() {
		nf.Stop()
	}()

	// Set up necessary variables
	var (
		recentRankLimit = 10
		sigch           = make(chan os.Signal, 1)
		ticker          = time.NewTicker(60 * time.Second)
		timeout         = time.NewTimer(3000 * time.Minute)
	)
	defer ticker.Stop()
	defer timeout.Stop()

	// Set up signal notification
	signal.Notify(sigch,
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT,
		syscall.SIGHUP, syscall.SIGUSR1, syscall.SIGUSR2,
	)

	// Context for managing goroutines
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle signals in a separate goroutine
	go handleSignals(sigch, cancel)

	// Process ranking in a separate goroutine
	go processRanking(ctx, nf, recentRankLimit, ticker)
	// Main event loop
	//for {
	select {
	case <-sigch:
		log.Println("Shutting down due to signal")
	case <-timeout.C:
		log.Println("Shutting down due to timeout")
	}
	log.Println("Exiting core function")
	//}
}

// Process ranking in a loop, handling errors and updating display
func processRanking(ctx context.Context, nf netflow.Interface, recentRankLimit int, ticker *time.Ticker) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			rank, err := nf.GetProcessRank(recentRankLimit, 60)
			if err != nil {
				log.Printf("GetProcessRank failed: %v", err)
				continue
			}

			clear()
			showTable(rank)
		}
	}
}

// Handle incoming signals
func handleSignals(sigch chan os.Signal, cancelFunc context.CancelFunc) {
	for sig := range sigch {
		log.Printf("Received signal: %s", sig)
		cancelFunc()
		return
	}
}

func showTable(ps []*netflow.Process) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"pid", "name", "exe", "inodes", "sum_in", "sum_out", "in_rate", "out_rate"})
	table.SetRowLine(true)

	items := [][]string{}
	urlProvider := rpc.UrlProvider{
		ServerEndpoint: "http://49.65.102.63:8334",
		CommonHeaders: rpc.CommonHeadersProvider{
			ImageVersion: "",
			DeviceId:     "c05803a7250ab9ccddb957122de312d0",
			BizType:      "2",
			Auth:         nil,
		},
	}
	client := rpc.CreateRpcClient(urlProvider)
	var in int64
	var out int64
	for _, po := range ps {
		inRate := humanBytes(po.TrafficStats.InRate)
		if po.TrafficStats.InRate > int64(constants.Thold) {
			inRate = constants.Red(inRate)
		}

		outRate := humanBytes(po.TrafficStats.OutRate)
		if po.TrafficStats.OutRate > int64(constants.Thold) {
			outRate = constants.Red(outRate)
		}

		item := []string{
			po.Pid,
			po.Name,
			po.Exe,
			cast.ToString(po.InodeCount),
			humanBytes(po.TrafficStats.In),
			humanBytes(po.TrafficStats.Out),
			inRate + "/s",
			outRate + "/s",
		}
		in += po.TrafficStats.InRate
		out += po.TrafficStats.OutRate
		items = append(items, item)
	}

	//遍历 如果又3个进程那么 累加所有3个进程流量的值
	now := time.Now().Unix()
	adjustedTime := now - (now % (1 * 60))
	down, _ := ConvertToMB(in)
	up, _ := ConvertToMB(out)
	testMonitorInfo := []rpc.MonitorInfo{
		{
			DownBandwidth: down,
			UpBandwidth:   up,
			CPUUsage:      0,
			DiskUsage:     0,
			MemUsage:      0,
			Timestamp:     adjustedTime,
		},
	}
	fmt.Printf("原始下载%d ,上传%d \n ", in, out)
	fmt.Printf("上报流量%v ,时间%s \n ", testMonitorInfo, time.Now().Format("2006-01-02 15:04:05"))
	err := client.ReportMonitorInfo(context.TODO(), testMonitorInfo)
	if err != nil {
		fmt.Print(err)
		return
	}

	items = append(items, []string{
		"total",
		"total",
		"total",
		cast.ToString(1),
		"",
		"",
		humanBytes(in) + "/s",
		humanBytes(out) + "/s",
	})
	table.AppendBulk(items)
	table.Render()
}

// Mock function to clear the console (replace with actual implementation)
func clear() {
	fmt.Printf("\x1b[2J")
}

// ConvertToMB 函数将 int64 字节值转换为 float64 兆字节值
func ConvertToMB(bytes int64) (float64, error) {
	// 将字节值转换为 MB（兆字节）
	mbStr := strconv.FormatFloat(float64(bytes)/1024/1024, 'f', 4, 64)
	// 将字符串转换为 float64
	mbFloat, err := strconv.ParseFloat(mbStr, 64)
	if err != nil {
		return 0, err // 如果转换错误，返回 0 和错误信息
	}
	return mbFloat, nil // 返回转换后的浮点数
}

func humanBytes(n int64) string {
	return humanize.Bytes(uint64(n))
}
