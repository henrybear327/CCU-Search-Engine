package main

import (
	"bytes"
	"container/list"
	"fmt"
	"log"
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
	robotsFile, statusCode := GetStaticSitePageSource(robotFileLink)
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

	fmt.Println("func", robotFileLink, robot.TestAgent("/", "CCU-assignment-bot"), robot.Sitemaps)
	manager.robot = robot
}

// ParseSiteMap extracts all links available
func (manager *Manager) parseSiteMap() {
	manager.urlQueue = list.New()
	if manager.robot == nil || len(manager.robot.Sitemaps) == 0 { // no sitemap
		manager.urlQueue.PushBack(manager.link)
	} else {
		// dfs
	}
}
