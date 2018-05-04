package main

import (
	"log"
	_ "net/http/pprof"
)

var conf config

func main() {
	parseConfigFile()

	// scheduler starts here!
	seedSiteList := getSeedSites()
	managers := prepareSeedSites(seedSiteList)
	log.Println("Manager count", len(managers))

	// run()
}
