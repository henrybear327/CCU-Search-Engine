package main

import (
	"fmt"
	_ "net/http/pprof"
)

var conf config

func main() {
	parseConfigFile()

	// scheduler starts here!
	seedSiteList := getSeedSites()
	managers := prepareSeedSites(seedSiteList)
	fmt.Println("Manager count", len(managers))
}
