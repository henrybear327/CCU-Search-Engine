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

	// var storage mongoDBStorage
	// storage.init()
	// // storage.ensureIndex("for every tld hub table", []string{"link"})
	// // storage.ensureIndex("for every tld page table", []string{"link"})
	// storage.deinit()
}
