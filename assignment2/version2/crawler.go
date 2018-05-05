package main

func startCrawling(managers map[string]*Manager) {
	dynamicLinkChannel := make(chan string)
	chromedpDone := make(chan bool)
	managerDone := make(chan bool)
	go getDynamicSitePageSource(dynamicLinkChannel, chromedpDone)
	for _, rec := range managers {
		dynamicLinkChannel <- rec.link

		for i := 0; i < conf.System.MaxGoRountinesPerSite; i++ {
			go rec.start(managerDone)
		}
	}

	for i := 0; i < conf.System.MaxGoRountinesPerSite*len(managers); i++ {
		<-managerDone
	}
	<-chromedpDone
}
