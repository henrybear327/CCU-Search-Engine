package main

import (
	_ "net/http/pprof"
)

var conf config

func main() {
	parseConfigFile()

	// scheduler starts here!
	seedSiteList := getSeedSites()
	managers := prepareSeedSites(seedSiteList)

	startCrawling(managers)
}
