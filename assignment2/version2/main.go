package main

import (
	_ "net/http/pprof"
)

var conf config

func main() {
	parseConfigFile()

	getDynamicSitePageSource("https://edition.cnn.com/")
	// getDynamicSitePageSource("https://www.npr.org/")

	// // scheduler starts here!
	// seedSiteList := getSeedSites()
	// managers := prepareSeedSites(seedSiteList)
	// log.Println("Manager count", len(managers))
}
