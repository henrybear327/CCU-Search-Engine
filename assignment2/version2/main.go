package main

import (
	_ "net/http/pprof"
)

func main() {
	var conf config
	parseConfigFile(&conf)

	// scheduler starts here!
	seedSiteList := getSeedSites(&conf)
	prepareSeedSites(seedSiteList)
}
