package main

import (
	"container/list"
	"crypto/sha1"
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/henrybear327/go-readability"
	"github.com/temoto/robotstxt"
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
	hub        map[string]int

	urlQueueLock   *sync.RWMutex
	urlFetchedLock *sync.RWMutex
	urlInQueueLock *sync.RWMutex
	hubLock        *sync.RWMutex

	distinctPagesFetched int
	useLinksFromXML      bool
	crawlDelay           time.Duration
	useStaticLoad        bool

	storage *mongoDBStorage
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

func (manager *Manager) isExternalSite(link string) bool {
	link = strings.TrimSpace(link)
	link = strings.ToLower(link)

	// fmt.Println(parsed.Host, parsed.Path)
	linkTLD := getTopLevelDomain(link)
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
func (manager *Manager) enqueue(link string, isPreprocessing bool) bool {
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

	if manager.isExternalSite(link) {
		return false
	}

	if manager.isBannedByRobotTXT(link) {
		return false
	}

	if isPreprocessing == false {
		if manager.isMultimediaFiles(link) {
			return false
		}
	}

	if manager.isInQueueOrFetched(link) {
		return false
	}

	manager.urlQueueLock.Lock()
	defer manager.urlQueueLock.Unlock()

	manager.urlInQueueLock.Lock()
	defer manager.urlInQueueLock.Unlock()

	if _, ok := manager.urlInQueue[link]; ok == false { // not in queue yet
		manager.urlQueue.PushBack(link)
		manager.urlInQueue[link] = true
		return true
	}

	return false
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

func (manager *Manager) doInDegreeCounting(link string) {
	if manager.isExternalSite(link) {
		return
	}

	if manager.isBannedByRobotTXT(link) {
		return
	}

	if manager.isMultimediaFiles(link) {
		return
	}

	manager.hubLock.Lock()
	defer manager.hubLock.Unlock()

	manager.hub[link]++
	manager.storage.hubUpsert(manager.tld, link, manager.hub[link])
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
		// dequeue -> fetch -> generate next link

		// dequeue
		nextLink, ok := manager.getNextURLFromQueue()
		if ok == false {
			break
		}

		// fetch
		fetchTime := time.Now()
		var pageSoruceForParsing []byte
		var titleForStoring string
		if manager.useStaticLoad {
			pageSource, statusCode := getStaticSitePageSource(nextLink)
			if statusCode != 200 {
				// TODO: disregard?
				time.Sleep(conf.System.minFetchTimeDuration)
				continue
			}
			pageSoruceForParsing = pageSource
			titleForStoring = getTitleFromPageSource(pageSource)
		} else {
			resultChannel := make(chan dynamicFetchingDataResult)
			query := dynamicFetchingDataQuery{link: nextLink, resultChannel: resultChannel}

			dynamicLinkChannel <- query
			result := <-resultChannel
			// fmt.Println(result.title, result.pageSource, result.requiresRestart)
			log.Println("dynamic query result", result.title, result.requiresRestart)
			// debugLog.Println("result", result.title, result.requiresRestart)

			if result.requiresRestart { // allocation error
				log.Println("put back", nextLink, "to queue front")
				manager.requeue(nextLink)
			} else {
				// TODO: title, pageSource integrity check...
				manager.addToFetched(nextLink)
			}
			pageSoruceForParsing = []byte(result.pageSource)
			titleForStoring = result.title
		}

		// generate next link
		log.Println("Parsed title", titleForStoring)
		nextLinkList := manager.getNextURLList(nextLink, pageSoruceForParsing)

		log.Println("Parsed links", len(nextLinkList))
		if len(nextLinkList) == 0 { // W.H.P this is multimedia!
			// ignore
		} else {
			if conf.Output.SavePageSource {
				saveHTMLFileFromString(getTopLevelDomain(nextLink), strings.Replace(nextLink[8:], "/", " ", -1)+".html", string(pageSoruceForParsing))
			}
			// log.Println("title", titleForStoring)

			test, err := readability.ParseFromPageSource(nextLink, pageSoruceForParsing, 5*time.Second)
			var mainText string
			var mainTextSHA1 string
			if err != nil {
				log.Println("failed.", err)
			} else {
				mainText = test.Content
				h := sha1.New()
				h.Write([]byte(mainText))
				bs := h.Sum(nil)
				mainTextSHA1 = fmt.Sprintf("%x\n", bs)
				// log.Println("main text", mainText)
			}
			manager.storage.sitePageUpsert(manager.tld, nextLink, fetchTime.Format(time.RFC3339), titleForStoring, string(pageSoruceForParsing), mainText, mainTextSHA1)

			for _, rec := range nextLinkList {
				// log.Println("Parsed link from", nextLink, rec)
				if manager.useLinksFromXML == false {
					ret := manager.enqueue(rec, false)
					if ret {
						// log.Println("Enqueue parsed link from", nextLink, rec)
					}
				}
				manager.doInDegreeCounting(rec)
			}
			log.Println("Queue size of", manager.tld, manager.urlQueue.Len())
		}

		if manager.useStaticLoad {
			time.Sleep(conf.System.minFetchTimeDuration)
		}
	}
	log.Println("Manager of", manager.link, "has finished")
}
