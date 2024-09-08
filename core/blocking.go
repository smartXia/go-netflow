package core

import (
	"context"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/rfyiamcool/go-netflow"
	"github.com/rfyiamcool/go-netflow/config"
	"github.com/rfyiamcool/go-netflow/constants"
	"github.com/rfyiamcool/go-netflow/rpc"
	"github.com/rfyiamcool/go-netflow/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var nf netflow.Interface

func Start(c config.Config) {
	var err error
	// Initialize netflow instance with error handling
	filter := ""
	if c.Filter != "" {
		filter = fmt.Sprintf("port %s", c.Filter)
	}
	nf, err = netflow.New(netflow.WithName(c.Nethogs), netflow.WithCaptureTimeout(12*30*24*60*time.Minute), netflow.WithPcapFilter(filter),
		netflow.WithQueueSize(2000000))
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
		recentRankLimit = 1
		sigch           = make(chan os.Signal, 1)
		ticker          = time.NewTicker(10 * time.Second)
		timeout         = time.NewTimer(365 * 30 * 24 * 60 * time.Minute)
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
	go processRanking(ctx, c, nf, recentRankLimit, ticker)
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

// 处理信号
func handleSignals(sigch chan os.Signal, cancelFunc context.CancelFunc) {
	for sig := range sigch {
		log.Printf("Received signal: %s", sig)
		cancelFunc()
		return
	}
}

// 在循环中处理排序，处理错误和更新显示
func processRanking(ctx context.Context, c config.Config, nf netflow.Interface, recentRankLimit int, ticker *time.Ticker) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			rank, err := nf.GetProcessRank(recentRankLimit, 5)
			if err != nil {
				log.Printf("GetProcessRank failed: %v", err)
				continue
			}
			showTable(c, rank)
			clear()
		}
	}
}

// 渲染和遍历
func showTable(c config.Config, ps []*netflow.Process) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"pid", "name", "exe", "inodes", "sum_in", "sum_out", "in_rate", "out_rate"})
	table.SetRowLine(true)
	var (
		items [][]string
		in    int64
		out   int64
	)
	for _, po := range ps {
		inRate, outRate := formatRates(po.TrafficStats.InRate, po.TrafficStats.OutRate)
		item := []string{po.Pid, po.Name, po.Exe, cast.ToString(po.InodeCount),
			utils.HumanBytes(po.TrafficStats.In * 8),
			utils.HumanBytes(po.TrafficStats.Out * 8),
			inRate,
			outRate,
		}
		//累加多进程级别的适配
		in += po.TrafficStats.InRate
		out += po.TrafficStats.OutRate
		items = append(items, item)
	}
	//上报流量信息
	reportHandler(in, out, c)
	table.AppendBulk(items)
	table.Render()
}

func reportHandler(in, out int64, c config.Config) {

	timestamp := time.Now().Unix()
	// 取整到最近的分钟
	adjustedTime := (timestamp / 60) * 60
	down, up := ConvertToMB(in), ConvertToMB(out)
	testMonitorInfo := []rpc.MonitorInfo{
		{
			DownBandwidth: down,
			UpBandwidth:   up,
			Timestamp:     adjustedTime,
		},
	}
	fmt.Printf("原始下载%d MiB ,上传%d MiB\n", in, out)
	fmt.Printf("上报流量%v ,时间%s\n", testMonitorInfo, time.Now().Format("2006-01-02 15:04:05"))
	urlProvider := rpc.UrlProvider{
		ServerEndpoint: c.Agent.ServerEndpoint,
		CommonHeaders: rpc.CommonHeadersProvider{
			ImageVersion: "",
			DeviceId:     c.Agent.DeviceId,
			BizType:      c.Agent.BizType,
			Ak:           c.Agent.AppKey,
			As:           c.Agent.AppSecret,
			Auth:         nil,
		},
	}
	client := rpc.CreateRpcClient(urlProvider)
	if err := client.ReportMonitorInfo(context.TODO(), testMonitorInfo); err != nil {
		fmt.Print(err)
		return
	}
}

func formatRates(inRate, outRate int64) (string, string) {
	// 使用 %.2f 来确保保留两位小数
	inRateStr := fmt.Sprintf("%.2f %s", float64(inRate*8)/1000/1000, "Mbit/s")
	if inRate > int64(constants.Thold) {
		inRateStr = constants.Red(inRateStr)
	}

	// 处理 outRate，保留两位小数
	outRateStr := fmt.Sprintf("%.2f %s", float64(outRate*8)/1000/1000, "Mbit/s")
	if outRate > int64(constants.Thold) {
		outRateStr = constants.Red(outRateStr)
	}

	return inRateStr, outRateStr
}

// ConvertToMB 函数将 int64 字节值转换为 float64 兆字节值
func ConvertToMB(bytes int64) float64 {
	// 将字节值转换为 MB（兆字节）
	mbStr := strconv.FormatFloat(float64(bytes)*8/1000/1000, 'f', 4, 64)
	// 将字符串转换为 float64
	mbFloat, _ := strconv.ParseFloat(mbStr, 64)
	return mbFloat // 返回转换后的浮点数
}

func clear() {
	fmt.Printf("\x1b[2J")
}
