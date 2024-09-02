package main

import (
	"flag"
	"github.com/rfyiamcool/go-netflow/core"
	log "github.com/sirupsen/logrus"
)

func main() {
	pname := flag.String("f", "", "choose p")
	flag.Parse()
	log.Info("core netflow sniffer")
	core.Start(*pname)
	//core(*pname)
	log.Info("netflow sniffer exit")
}
