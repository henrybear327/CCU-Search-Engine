package main

import (
	"log"
	"time"
)

func startCrawling(managers map[string]*Manager) {
	// init channels
	dynamicLinkChannel := make(chan dynamicFetchingDataQuery)
	managerDone := make(chan bool)

	// init storage
	var storage mongoDBStorage
	storage.init()
	defer storage.deinit()
	storage.ensureIndex("hub", "tld", "link")
	storage.ensureIndex("sitePage", "tld", "link")

	// think of creating a daemon
	// for creating it, we make channels
	// for using it, we use channels
	go getDynamicSitePageSource(dynamicLinkChannel)
	for i := 0; i < conf.System.MaxGoRountinesPerSite; i++ {
		for _, rec := range managers {
			go rec.start(managerDone, dynamicLinkChannel, &storage)
		}
		// if no sitemap.xml, only one thread will be alive since queue size = 1 can only serve 1 thread QQ (e.g. npr.org)
		// reduce system load
		if i == 0 && conf.System.MaxGoRountinesPerSite > 1 {
			log.Println("Delay lanuching...")
			time.Sleep(conf.System.minFetchTimeDuration * 3)
			log.Println("Launching all goroutines")
		}
	}

	globalTimeout := time.After(conf.System.maxRunningTimeDuration)
	log.Println("Global timeout", globalTimeout)
	for i := 0; i < conf.System.MaxGoRountinesPerSite*len(managers); i++ {
		select {
		case <-managerDone:
		case <-globalTimeout:
			log.Println("Global timeout! Ending!")
			return
		}
	}
}
