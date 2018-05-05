package main

func startCrawling(managers map[string]*Manager) {
	dynamicLinkChannel := make(chan dynamicFetchingDataQuery)
	managerDone := make(chan bool)

	// think of creating a daemon
	// for creating it, we make channels
	// for using it, we use channels
	go getDynamicSitePageSource(dynamicLinkChannel)
	for _, rec := range managers {
		for i := 0; i < conf.System.MaxGoRountinesPerSite; i++ {
			go rec.start(managerDone, dynamicLinkChannel)
		}
	}

	for i := 0; i < conf.System.MaxGoRountinesPerSite*len(managers); i++ {
		<-managerDone
	}
}
