package main

import (
	"bytes"
	"compress/gzip"
	"container/list"
	"encoding/xml"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/temoto/robotstxt"
)

type sitemapSection struct {
	Sitemap []sitemapSectionNode `xml:"sitemap"`
}

type sitemapSectionNode struct {
	Loc string `xml:"loc"` // url
}

type sitemapURL struct {
	URL []sitemapURLNode `xml:"url"`
}

type sitemapURLNode struct {
	Loc string `xml:"loc"` // url
}

func (manager *Manager) processgzFile(link string) {
	// log.Println("processgzFile")
	compressedPageSource, statusCode := getStaticSitePageSource(link)
	if statusCode != 200 {
		log.Println("[error] gz file download err")
		return
	}

	manager.addToFetched(link) // Crucial! Techbang infinit loop fucking bug

	res := bytes.NewReader(compressedPageSource)

	gzf, err := gzip.NewReader(res)
	if err != nil {
		log.Println("[error] gzip new reader", err)
	}

	pageSource, err := ioutil.ReadAll(gzf)
	if err != nil {
		log.Println("[error] gzip readAll", err)
	}

	manager.parseXMLContent(pageSource)
}

func (manager *Manager) parseXMLContent(pageSource []byte) {
	// log.Println("parseXMLContent")
	// parse (section xml)
	var data sitemapSection
	err := xml.Unmarshal(pageSource, &data)
	if err != nil {
		log.Println("[error] parse sitemapSection error", err)
		return
	}

	if len(data.Sitemap) == 0 {
		// not section xml, it's link xml!
		var data sitemapURL
		err := xml.Unmarshal(pageSource, &data)
		if err != nil {
			color.Set(color.FgRed)
			log.Println("[error] parse sitemapURL error", err)
			color.Unset()
			return
		}

		blocking := make(chan bool)
		for _, rec := range data.URL {
			// fmt.Println("link", rec.Loc)
			go manager.generateLinksFromSitemap(rec.Loc, blocking)
			<-blocking // no need to async since it's link!
		}
	} else {
		blocking := make(chan bool)
		for _, rec := range data.Sitemap {
			// fmt.Println("link", rec.Loc)
			go manager.generateLinksFromSitemap(rec.Loc, blocking) // multi process xml files
			// <-blocking
		}

		for i := 0; i < len(data.Sitemap); i++ {
			<-blocking
		}
	}
}

func (manager *Manager) generateLinksFromSitemap(link string, done chan bool) {
	link = strings.TrimSpace(link)
	// log.Println("Now DFS", link, manager.isInQueueOrFetched(link))
	if manager.isInQueueOrFetched(link) {
		done <- true
		return
	}

	if manager.isExternalSite(link) {
		// fmt.Println("Is external site", link)
		done <- true
		return
	}

	if strings.HasSuffix(link, ".gz") {
		manager.processgzFile(link)
		done <- true
		return
	}

	if strings.HasSuffix(link, ".xml") == false {
		manager.enqueue(link, true)

		done <- true
		return
	}

	pageSource, statusCode := getStaticSitePageSource(link)
	if statusCode != 200 {
		done <- true
		return
	}

	manager.addToFetched(link)
	manager.parseXMLContent(pageSource)
	done <- true
}

// ParseSiteMap extracts all links available
func (manager *Manager) parseSitemap() {
	startParsing := time.Now()

	manager.urlQueue = list.New()

	if manager.robot == nil || len(manager.robot.Sitemaps) == 0 { // no sitemap
		manager.enqueue(manager.link, true)
	} else {
		manager.useLinksFromXML = true

		done := make(chan bool)
		for _, mapLink := range manager.robot.Sitemaps {
			// dfs
			go manager.generateLinksFromSitemap(mapLink, done)
		}

		for i := 0; i < len(manager.robot.Sitemaps); i++ {
			<-done
		}
	}

	if manager.urlQueue.Len() == 0 {
		manager.enqueue(manager.link, true)
	}

	// for e := manager.urlQueue.Front(); e != nil; e = e.Next() {
	// 	fmt.Println(e.Value.(string))
	// }

	elapsedParsing := time.Since(startParsing)
	log.Println(manager.link, "Initial queue size", manager.urlQueue.Len(), "Parsing XML takes", elapsedParsing)
}

// ParseRobotsTxt attempts parses the robots.txt file of the given link
func (manager *Manager) parseRobotsTxt() {
	robotFileLink := manager.link + "/robots.txt"
	robotsFile, statusCode := getStaticSitePageSource(robotFileLink)
	if statusCode != 200 {
		color.Set(color.FgRed)
		log.Println("[error] error fetching robots.txt for site", robotFileLink, statusCode)
		color.Unset()

		manager.robot = nil
		return
	}

	robot, err := robotstxt.FromStatusAndBytes(statusCode, robotsFile)
	if err != nil {
		color.Set(color.FgRed)
		log.Println("[error] Error parsing robots.txt", err)
		color.Unset()
	}

	manager.robot = robot
	// fmt.Println("parseRobotsTxt", robotFileLink, manager.robot.TestAgent("/", "CCU-assignment-bot"), manager.robot.Sitemaps)
}

func (manager *Manager) preprocess(done chan bool) {
	// parse robots.txt
	manager.parseRobotsTxt()

	// parse sitemap.xml (prepare queue)
	manager.parseSitemap()

	done <- true
}
