package main

func startCrawling(managers map[string]*Manager) {
	dynamicLinkChannel := make(chan string)
	done := make(chan bool)
	go getDynamicSitePageSource(dynamicLinkChannel, done)
	for _, rec := range managers {
		dynamicLinkChannel <- rec.link
	}
	<-done
}
