package main

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
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

func isInvalidSuffix(link string) bool {
	if link == "void(0)" {
		return true
	}
	if strings.HasSuffix(link, "mailto://") {
		return true
	}
	if strings.HasSuffix(link, "javascript://") {
		return true
	}

	return false
}

func isValidURL(link string) bool {
	link = strings.ToLower(strings.TrimSpace(link))

	if isInvalidSuffix(link) {
		return false
	}

	return true
}

func (manager *Manager) generateNextURLList(pageSource []byte) []string {
	nextURLs := []string{}

	// Load the HTML document
	res := bytes.NewReader(pageSource)
	doc, err := goquery.NewDocumentFromReader(res)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		// 1. must have href
		// 2. concat url
		// 3. enqueue
		// 4. hub counting
		fmt.Println(i, s, s.Text())
	})

	return nextURLs
}
