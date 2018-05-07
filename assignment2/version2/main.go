package main

import (
	_ "net/http/pprof"
)

var conf config

func main() {
	parseConfigFile()

	// scheduler starts here!
	seedSiteList, seedSiteOption := getSeedSites()
	managers := prepareSeedSites(seedSiteList, seedSiteOption)

	startCrawling(managers)
}
