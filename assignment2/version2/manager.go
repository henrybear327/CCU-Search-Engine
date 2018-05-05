package main

import (
	"container/list"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/temoto/robotstxt"
	"golang.org/x/net/publicsuffix"
)

// Manager is the heart of every seed website
// Be aware of locking
type Manager struct {
	link       string
	tld        string
	robot      *robotstxt.RobotsData
	urlQueue   *list.List
	urlFetched map[string]bool
	urlInQueue map[string]bool

	urlQueueLock   *sync.RWMutex
	urlFetchedLock *sync.RWMutex
	urlInQueueLock *sync.RWMutex

	distinctPagesFetched int
	useLinksFromXML      bool
	crawlDelay           time.Duration
	useStaticLoad        bool
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

func getTLD(link string) string {
	linkTLD, err := publicsuffix.EffectiveTLDPlusOne(link)
	if err != nil {
		log.Println("isExternalSite EffectiveTLDPlusOne err", err)
		return ""
	}
	return linkTLD
}

func (manager *Manager) isExternalSite(link string) bool {
	link = strings.TrimSpace(link)
	link = strings.ToLower(link)

	parsed, err := url.Parse(link)
	if err != nil {
		log.Println("isExternalSite url parse err", err)
		return true // can't parse, disregard
	}

	// fmt.Println(parsed.Host, parsed.Path)
	linkTLD := getTLD(parsed.Host)
	// if manager.tld != linkTLD {
	// 	fmt.Println("cmp isExternalSite", manager.tld, link, linkTLD)
	// }

	return manager.tld != linkTLD
}

// call this afer fetching
func (manager *Manager) addToFetched(link string) {
	link = strings.TrimSpace(link)

	// if in InQueue map, move it from there to fetched
	manager.urlInQueueLock.Lock()
	defer manager.urlInQueueLock.Unlock()
	if _, ok := manager.urlInQueue[link]; ok == true {
		// add to fetched set
		manager.urlFetchedLock.Lock()
		defer manager.urlFetchedLock.Unlock()
		manager.urlFetched[link] = true

		// remove from InQueue
		delete(manager.urlInQueue, link)

		manager.distinctPagesFetched++
		if manager.distinctPagesFetched >= conf.System.MaxDistinctPagesToFetchPerSite {
			// TODO: end go routine
		}
		return
	}

	// somehow it's fetched twice?
	if manager.isInQueueOrFetched(link) {
		log.Println("Weird shit, fetched twice?")
		return
	}
}

func (manager *Manager) isMultimediaFiles(link string) bool {
	// TODO: maybe use (html|php|...) match?
	// https://fileinfo.com/filetypes/common
	regex := "^.*\\.(doc|docx|odt|csv|ppt|pptx|wav|wma|jpg|png|gif|jpeg|mp3|mp4|mov|avi|flv)$"
	matched, err := regexp.MatchString(regex, link)
	if err != nil {
		log.Println("isMultimediaFiles", err)
		return true
	}

	if matched {
		log.Println("eliminated", link)
	}
	return matched
}

func (manager *Manager) isBannedByRobotTXT(link string) bool {
	link = strings.ToLower(strings.TrimSpace(link))

	if manager.robot != nil && manager.robot.TestAgent(link, "CCU-Assignment-Bot") == false {
		return true
	}
	return false
}

// call this to queue url
func (manager *Manager) enqueue(link string, isPreprocessing bool) {
	link = strings.TrimSpace(link)
	/*
		Disgard link if
		1. v already in queue
		2. v already fetched
		3. v ending with unwanted filetype
		4. v link is going to external site
		5. v against robot rules
		6. x main text hash collision (?)
	*/

	if manager.isInQueueOrFetched(link) {
		return
	}

	if manager.isExternalSite(link) {
		return
	}

	if manager.isBannedByRobotTXT(link) {
		return
	}

	if isPreprocessing == false {
		if manager.isMultimediaFiles(link) {
			return
		}
	}

	manager.urlQueueLock.Lock()
	defer manager.urlQueueLock.Unlock()

	manager.urlInQueueLock.Lock()
	defer manager.urlInQueueLock.Unlock()

	if _, ok := manager.urlInQueue[link]; ok == false { // not in queue yet
		manager.urlQueue.PushBack(link)
		manager.urlInQueue[link] = true
	}
}

// for restart
func (manager *Manager) requeue(link string) {
	manager.urlQueueLock.Lock()
	defer manager.urlQueueLock.Unlock()

	manager.urlQueue.PushFront(link)
}

func (manager *Manager) getNextURLFromQueue() string {
	manager.urlQueueLock.Lock()
	defer manager.urlQueueLock.Unlock()

	// dequeue, but don't remove it from InQueue map
	nextLink := manager.urlQueue.Front()
	manager.urlQueue.Remove(nextLink)

	return (*nextLink).Value.(string)
}

func (manager *Manager) hasNextURL() bool {
	manager.urlQueueLock.Lock()
	defer manager.urlQueueLock.Unlock()

	return manager.urlQueue.Len() > 0
}

func (manager *Manager) start(done chan bool, dynamicLinkChannel chan dynamicFetchingDataQuery) {
	defer func(done chan bool) {
		done <- true
	}(done)

	log.Println("Manager of ", manager.link, "is started")
	for manager.hasNextURL() {
		// var title, pageSource string
		if manager.useLinksFromXML {
			// simply dequeue and fetch
		} else {
			nextLink := manager.getNextURLFromQueue()

			// dequeue -> fetch page -> generate next links
			resultChannel := make(chan dynamicFetchingDataResult)
			query := dynamicFetchingDataQuery{link: nextLink, resultChannel: resultChannel}

			dynamicLinkChannel <- query
			result := <-resultChannel
			// fmt.Println(result.title, result.pageSource, result.requiresRestart)
			fmt.Println(result.title, result.requiresRestart)

			if result.requiresRestart {
				manager.requeue(nextLink)
			} else {
				// TODO: title, pageSource integrity check...
				manager.addToFetched(nextLink)
			}
		}
	}
	log.Println("Manager of ", manager.link, "has finished")
}
