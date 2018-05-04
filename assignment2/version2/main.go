package main

import (
	_ "net/http/pprof"
)

var conf config

func main() {
	parseConfigFile()

	done := make(chan bool)
	go getDynamicSitePageSource("https://edition.cnn.com/", done)
	<-done
	go getDynamicSitePageSource("https://www.npr.org/", done)
	<-done

	// // scheduler starts here!
	// seedSiteList := getSeedSites()
	// managers := prepareSeedSites(seedSiteList)
	// log.Println("Manager count", len(managers))
}
