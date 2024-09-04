package main

import (
	"flag"
	"github.com/rfyiamcool/go-netflow/config"
	"github.com/rfyiamcool/go-netflow/core"
	log "github.com/sirupsen/logrus"
)

func main() {
	pname := flag.String("f", "", "choose p")
	configPathPtr := flag.String("config", "/usr/local/super-agent/config.yaml", "super-agent config file path")
	flag.Parse()
	log.Info("core netflow sniffer")
	configSuperAgent := config.GetConfig(*configPathPtr)
	configSuperAgent.Nethogs = *pname
	if *configPathPtr == "" {
		panic("super-agent config path must be provided")
	}
	core.Start(configSuperAgent)
	//core(*pname)
	log.Info("netflow sniffer exit")
}
