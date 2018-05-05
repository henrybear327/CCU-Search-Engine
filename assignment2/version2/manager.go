package main

import (
	"container/list"
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
// Don't forget to acquire lock
// Be aware of deadlock
// urlFetchedLock -> urlQueueLock -> urlInQueueLock
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
	manager.urlFetchedLock.RLock()
	defer manager.urlFetchedLock.RUnlock()
	manager.urlInQueueLock.RLock()
	defer manager.urlInQueueLock.RUnlock()

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
	// log.Println("addToFetched", link)
	link = strings.TrimSpace(link)

	// add to fetched set
	manager.urlFetchedLock.Lock()
	defer manager.urlFetchedLock.Unlock()
	if _, ok := manager.urlFetched[link]; ok == true {
		return
	}
	manager.urlFetched[link] = true

	manager.distinctPagesFetched++
	if manager.distinctPagesFetched >= conf.System.MaxDistinctPagesToFetchPerSite {
		// TODO: end go routine
	}

	// if in InQueue map, move it from there to fetched
	manager.urlInQueueLock.Lock()
	defer manager.urlInQueueLock.Unlock()
	if _, ok := manager.urlInQueue[link]; ok == true {
		// remove from InQueue
		delete(manager.urlInQueue, link)
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

func (manager *Manager) getNextURLFromQueue() (string, bool) {
	manager.urlQueueLock.Lock()
	defer manager.urlQueueLock.Unlock()

	// dequeue, but don't remove it from InQueue map
	if manager.urlQueue.Len() > 0 {
		nextLink := manager.urlQueue.Front()
		manager.urlQueue.Remove(nextLink)

		return (*nextLink).Value.(string), true
	}

	return "", false
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

	// fileName := manager.link[9:13] + ".log"
	// logFile, err := os.Create(fileName)
	// defer logFile.Close()
	// if err != nil {
	// 	log.Fatalln("manager log", err)
	// }

	// debugLog := log.New(logFile, manager.link[9:], log.LstdFlags)
	// debugLog.SetFlags(debugLog.Flags() | log.LstdFlags)

	for manager.hasNextURL() {
		// dequeue -> fetch
		nextLink, ok := manager.getNextURLFromQueue()
		if ok == false {
			break
		}

		resultChannel := make(chan dynamicFetchingDataResult)
		query := dynamicFetchingDataQuery{link: nextLink, resultChannel: resultChannel}

		dynamicLinkChannel <- query
		result := <-resultChannel
		// fmt.Println(result.title, result.pageSource, result.requiresRestart)
		log.Println("dynamic query result", result.title, result.requiresRestart)
		// debugLog.Println("result", result.title, result.requiresRestart)

		if result.requiresRestart {
			log.Println("put back", nextLink, "to queue front")
			manager.requeue(nextLink)
		} else {
			// TODO: title, pageSource integrity check...
			manager.addToFetched(nextLink)
		}

		// var title, pageSource string
		if manager.useLinksFromXML {
			// nothing more to do
		} else {
			// TODO: generate next links
		}
	}
	log.Println("Manager of", manager.link, "has finished")
}
