package main

import (
	_ "net/http/pprof"
)

var conf config

func main() {
	parseConfigFile()

	link := make(chan string)
	done := make(chan bool)
	go getDynamicSitePageSource(link, done)
	links := []string{"https://edition.cnn.com/", "https://www.npr.org/", "https://www.techbang.com", "https://google.com", "http://www.ccu.edu.tw", "https://bbc.co", "https://youtube.com"}
	for _, rec := range links {
		link <- rec
	}
	<-done

	// // scheduler starts here!
	// seedSiteList := getSeedSites()
	// managers := prepareSeedSites(seedSiteList)
	// log.Println("Manager count", len(managers))
}
