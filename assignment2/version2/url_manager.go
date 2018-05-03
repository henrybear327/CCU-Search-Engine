package main

import (
	"container/list"
	"sync"

	"github.com/temoto/robotstxt"
)

// Manager is the heart of every seed website
type Manager struct {
	link       string
	robot      *robotstxt.RobotsData
	urlQueue   *list.List
	urlFetched map[string]bool
	urlInQueue map[string]bool
	conf       *config

	urlQueueLock *sync.Mutex

	distinctPagesFetched int
}

func (manager *Manager) preprocess(done chan bool) {
	// parse robots.txt
	manager.parseRobotsTxt()

	// parse sitemap.xml (prepare queue)
	manager.parseSiteMap()

	done <- true
}

func (manager *Manager) enqueue(link string) {
	/*
		Disgard link if
		1. already in queue
		2. already fetched
		3. main text hash collision (?)
		4. ending with unwanted filetype
		5. going external
	*/

	manager.urlQueueLock.Lock()
	defer manager.urlQueueLock.Unlock()

	manager.urlQueue.PushBack(link)
	manager.distinctPagesFetched++

	if manager.distinctPagesFetched >= manager.conf.System.MaxDistinctPagesToFetchPerSite {
		// TODO: end go routine
	}
}
