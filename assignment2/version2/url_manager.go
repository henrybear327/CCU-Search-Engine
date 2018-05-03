package main

import (
	"container/list"
	"strings"
	"sync"

	"github.com/temoto/robotstxt"
)

// Manager is the heart of every seed website
// Be aware of locking
type Manager struct {
	link       string
	robot      *robotstxt.RobotsData
	urlQueue   *list.List
	urlFetched map[string]bool
	urlInQueue map[string]bool
	conf       *config

	urlQueueLock   *sync.RWMutex
	urlFetchedLock *sync.RWMutex
	urlInQueueLock *sync.RWMutex

	distinctPagesFetched int
}

func (manager *Manager) preprocess(done chan bool) {
	// parse robots.txt
	manager.parseRobotsTxt()

	// parse sitemap.xml (prepare queue)
	manager.parseSiteMap()

	done <- true
}

func (manager *Manager) isInQueueOrFetched(link string) bool {
	manager.urlInQueueLock.RLock()
	defer manager.urlInQueueLock.RUnlock()
	manager.urlFetchedLock.RLock()
	defer manager.urlFetchedLock.RUnlock()

	if _, ok := manager.urlInQueue[link]; ok {
		return true
	}
	if _, ok := manager.urlFetched[link]; ok {
		return true
	}

	return false
}

func (manager *Manager) enqueue(link string) {
	link = strings.TrimSpace(link)
	/*
		Disgard link if
		1. v already in queue
		2. v already fetched
		3. x main text hash collision (?)
		4. x ending with unwanted filetype
		5. x link is going to external site
	*/

	if manager.isInQueueOrFetched(link) {
		return
	}

	manager.urlQueueLock.Lock()
	defer manager.urlQueueLock.Unlock()

	manager.urlInQueueLock.Lock()
	defer manager.urlInQueueLock.Unlock()

	if _, ok := manager.urlInQueue[link]; ok == false { // not in queue yet
		manager.urlQueue.PushBack(link)
		manager.distinctPagesFetched++
		manager.urlInQueue[link] = true

		if manager.distinctPagesFetched >= manager.conf.System.MaxDistinctPagesToFetchPerSite {
			// TODO: end go routine
		}
	}
}
