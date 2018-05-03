package main

import (
	"bytes"
	"compress/gzip"
	"container/list"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
	"github.com/temoto/robotstxt"
)

// ParseAlexaTopSites is a function that takes html source code of Alexa top 50 site
// and parse the top sites out to an array of strings
func parseAlexaTopSites(pageSource []byte) []string {
	// Load the HTML document
	startParsing := time.Now()

	res := bytes.NewReader(pageSource)
	doc, err := goquery.NewDocumentFromReader(res)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	topLinkList := make([]string, 0)
	// log.Println("Top Alexa sites are:")
	// #alx-content > div > div > section.page-product-content.summary > span > span > div > div > div.listings.table > div:nth-child(2) > div.td.DescriptionCell > p > a
	doc.Find("div.td.DescriptionCell p").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		link := s.Find("a").Text()
		// log.Println(link)

		topLinkList = append(topLinkList, link)
		// topLinkList = append(topLinkList, "https://www."+link)
	})

	elapsedParsing := time.Since(startParsing)
	log.Printf("Parsing top Alexa sites took %s", elapsedParsing)

	return topLinkList
}

// ParseRobotsTxt attempts parses the robots.txt file of the given link
func (manager *Manager) parseRobotsTxt() {
	robotFileLink := manager.link + "/robots.txt"
	robotsFile, statusCode := getStaticSitePageSource(robotFileLink)
	if statusCode != 200 {
		color.Set(color.FgRed)
		log.Println("Error fetching robots.txt for site", robotFileLink, statusCode)
		color.Unset()

		manager.robot = nil
		return
	}

	robot, err := robotstxt.FromStatusAndBytes(statusCode, robotsFile)
	if err != nil {
		color.Set(color.FgRed)
		log.Println("Error parsing robots.txt", err)
		color.Unset()
	}

	// fmt.Println("func", robotFileLink, robot.TestAgent("/", "CCU-assignment-bot"), robot.Sitemaps)
	manager.robot = robot
}

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
	// fmt.Println("processgzFile")
	compressedPageSource, statusCode := getStaticSitePageSource(link)
	if statusCode != 200 {
		return
	}

	manager.addToFetched(link) // Crucial! Techbang infinit loop fucking bug

	res := bytes.NewReader(compressedPageSource)

	gzf, err := gzip.NewReader(res)
	if err != nil {
		log.Println(err)
	}

	pageSource, err := ioutil.ReadAll(gzf)
	if err != nil {
		log.Println(err)
	}

	manager.parseXMLContent(pageSource)
}

func (manager *Manager) parseXMLContent(pageSource []byte) {
	// fmt.Println("parseXMLContent")
	// parse (section xml)
	var data sitemapSection
	err := xml.Unmarshal(pageSource, &data)
	if err != nil {
		log.Println("parse sitemapSection error", err)
		return
	}

	if len(data.Sitemap) == 0 {
		// not section xml, it's link xml!
		var data sitemapURL
		err := xml.Unmarshal(pageSource, &data)
		if err != nil {
			color.Set(color.FgRed)
			log.Println("parse sitemapURL error", err)
			color.Unset()
			return
		}

		blocking := make(chan bool)
		for _, rec := range data.URL {
			// fmt.Println("link", rec.Loc)
			go manager.generateLinksFromSitemap(rec.Loc, blocking)
			// <-blocking
		}

		for i := 0; i < len(data.URL); i++ {
			<-blocking
		}
	} else {
		blocking := make(chan bool)
		for _, rec := range data.Sitemap {
			// fmt.Println("link", rec.Loc)
			go manager.generateLinksFromSitemap(rec.Loc, blocking)
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
		fmt.Println("Is external site", link)
		done <- true
		return
	}

	if strings.HasSuffix(link, ".gz") {
		manager.processgzFile(link)
		done <- true
		return
	}

	if strings.HasSuffix(link, ".xml") == false {
		manager.enqueue(link)

		done <- true
		return
	}

	pageSource, statusCode := getStaticSitePageSource(link)
	if statusCode != 200 {
		done <- true
		return
	}

	manager.parseXMLContent(pageSource)
	done <- true
}

// ParseSiteMap extracts all links available
func (manager *Manager) parseSiteMap() {
	startParsing := time.Now()

	manager.urlQueue = list.New()

	if manager.robot == nil || len(manager.robot.Sitemaps) == 0 { // no sitemap
		manager.enqueue(manager.link)
	} else {
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
		manager.enqueue(manager.link)
	}

	// for e := manager.urlQueue.Front(); e != nil; e = e.Next() {
	// 	fmt.Println(e.Value.(string))
	// }

	elapsedParsing := time.Since(startParsing)
	fmt.Println(manager.link, "initial queue size", manager.urlQueue.Len())
	fmt.Println("Parsing XML takes", elapsedParsing)
}
