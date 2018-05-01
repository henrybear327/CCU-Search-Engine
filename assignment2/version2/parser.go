package main

import (
	"bytes"
	"log"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// ParseAlexaTopSites is a function that takes a html source code string of Alexa top 50 site
// and parse the top sites
func ParseAlexaTopSites(pageSource []byte) []string {
	// Load the HTML document
	startParsing := time.Now()

	res := bytes.NewReader(pageSource)
	doc, err := goquery.NewDocumentFromReader(res)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	topURLList := make([]string, 0)
	log.Println("Top Alexa sites are:")
	doc.Find("div.tr.site-listing").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		url := s.Find("a").Text()
		log.Println(url)

		topURLList = append(topURLList, url)
	})

	elapsedParsing := time.Since(startParsing)
	log.Printf("Parsing top Alexa sites took %s", elapsedParsing)
	return topURLList
}
