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

	dynamicLinkChannel := make(chan string)
	done := make(chan bool)
	go getDynamicSitePageSource(dynamicLinkChannel, done)
	// links := []string{"https://edition.cnn.com/", "https://www.npr.org/", "https://www.techbang.com", "https://google.com", "http://www.ccu.edu.tw", "https://bbc.co", "https://youtube.com"}
	// for _, rec := range links {
	// 	dynamicLinkChannel <- rec
	// }
	for _, rec := range managers {
		dynamicLinkChannel <- rec.link
	}
	<-done
}
