package main

import (
	"container/list"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"
	"sync"

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
		fmt.Println("isExternalSite EffectiveTLDPlusOne err", err)
		return ""
	}
	return linkTLD
}

func (manager *Manager) isExternalSite(link string) bool {
	link = strings.TrimSpace(link)
	link = strings.ToLower(link)

	parsed, err := url.Parse(link)
	if err != nil {
		fmt.Println("isExternalSite url parse err", err)
		return true // can't parse, disregard
	}

	// fmt.Println(parsed.Host, parsed.Path)
	linkTLD := getTLD(parsed.Host)
	// if manager.tld != linkTLD {
	// 	fmt.Println("cmp isExternalSite", manager.tld, link, linkTLD)
	// }

	return manager.tld != linkTLD
}

func (manager *Manager) addToFetched(link string) {
	link = strings.TrimSpace(link)

	if manager.isInQueueOrFetched(link) {
		return
	}

	manager.urlFetchedLock.Lock()
	defer manager.urlFetchedLock.Unlock()

	manager.urlFetched[link] = true

	manager.distinctPagesFetched++
	if manager.distinctPagesFetched >= conf.System.MaxDistinctPagesToFetchPerSite {
		// TODO: end go routine
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

func (manager *Manager) enqueue(link string, isPreprocessing bool) {
	link = strings.TrimSpace(link)
	/*
		Disgard link if
		1. v already in queue
		2. v already fetched
		3. x main text hash collision (?)
		4. x ending with unwanted filetype
		5. v link is going to external site
		6. v against robot rules
	*/

	if manager.isInQueueOrFetched(link) {
		return
	}

	if manager.robot != nil && manager.robot.TestAgent(link, "CCU-Assignment-Bot") == false {
		return
	}

	if manager.isExternalSite(link) {
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
