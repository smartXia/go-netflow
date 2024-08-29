package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/rfyiamcool/go-netflow/rpc"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"

	"github.com/rfyiamcool/go-netflow"
)

var (
	nf netflow.Interface

	yellow  = color.New(color.FgYellow).SprintFunc()
	red     = color.New(color.FgRed).SprintFunc()
	info    = color.New(color.FgGreen).SprintFunc()
	blue    = color.New(color.FgBlue).SprintFunc()
	magenta = color.New(color.FgHiMagenta).SprintFunc()
)

func start(pname string) {
	var err error

	nf, err = netflow.New(netflow.WithName(pname))
	if err != nil {
		log.Fatal(err)
	}

	err = nf.Start()
	if err != nil {
		log.Fatal(err)
	}

	var (
		recentRankLimit = 10

		sigch   = make(chan os.Signal, 1)
		ticker  = time.NewTicker(3 * time.Second)
		timeout = time.NewTimer(3000 * time.Second)
	)

	signal.Notify(sigch,
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT,
		syscall.SIGHUP, syscall.SIGUSR1, syscall.SIGUSR2,
	)

	defer func() {
		nf.Stop()
	}()

	go func() {
		for {
			<-ticker.C
			rank, err := nf.GetProcessRank(recentRankLimit, 3)
			if err != nil {
				log.Errorf("GetProcessRank failed, err: %s", err.Error())
				continue
			}

			clear()
			showTable(rank)
		}
	}()

	for {
		select {
		case <-sigch:
			return

		case <-timeout.C:
			return
		}
	}
}

func stop() {
	if nf == nil {
		return
	}

	nf.Stop()
}

const thold = 1024 * 1024 // 1mb

func clear() {
	fmt.Printf("\x1b[2J")
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

	for _, po := range ps {
		inRate := humanBytes(po.TrafficStats.InRate)
		if po.TrafficStats.InRate > int64(thold) {
			inRate = red(inRate)
		}

		outRate := humanBytes(po.TrafficStats.OutRate)
		if po.TrafficStats.OutRate > int64(thold) {
			outRate = red(outRate)
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

		items = append(items, item)
		now := time.Now().Unix()
		adjustedTime := now - (now % (0.5 * 60))
		testMonitorInfo := []rpc.MonitorInfo{
			{
				DownBandwidth: float64(po.TrafficStats.In),
				UpBandwidth:   float64(po.TrafficStats.Out),
				CPUUsage:      0,
				DiskUsage:     0,
				MemUsage:      0,
				Timestamp:     adjustedTime,
			},
		}
		fmt.Printf("上报流量%v \n", testMonitorInfo)
		err := client.ReportMonitorInfo(context.TODO(), testMonitorInfo)
		if err != nil {
			fmt.Print(err)
			return
		}
	}

	table.AppendBulk(items)
	table.Render()

}
func getFloat(str string) float64 {
	floatValue, _ := strconv.ParseFloat(str, 64)
	return floatValue
}
func humanBytes(n int64) string {
	return humanize.Bytes(uint64(n))
}

func main() {
	pname := flag.String("f", "", "choose p")
	flag.Parse()
	log.Info("start netflow sniffer")

	start(*pname)
	stop()

	log.Info("netflow sniffer exit")
}
