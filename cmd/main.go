package main

import (
	"flag"
	"github.com/rfyiamcool/go-netflow/config"
	"github.com/rfyiamcool/go-netflow/core"
	log "github.com/sirupsen/logrus"
	"net/http"
	_ "net/http/pprof" // 引入 pprof 包
)

func main() {
	pname := flag.String("f", "", "choose p")
	filter := flag.String("p", "", "choose port")
	configPathPtr := flag.String("config", "/usr/local/super-agent/config.yaml", "super-agent config file path")
	flag.Parse()

	// 启动 pprof 服务
	go func() {
		log.Info("Starting pprof server on :6060")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Fatal("pprof server failed:", err)
		}
	}()

	log.Info("core netflow sniffer")
	configSuperAgent := config.GetConfig(*configPathPtr)
	configSuperAgent.Nethogs = *pname
	configSuperAgent.Filter = *filter
	if *configPathPtr == "" {
		panic("super-agent config path must be provided")
	}

	core.Start(configSuperAgent)
	//core2.StartOut(configSuperAgent)
	//core(*pname)
	log.Info("netflow sniffer exit")
}
