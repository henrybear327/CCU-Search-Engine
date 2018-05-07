package main

import (
	"bytes"
	"log"
	"net/url"
	"regexp"
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

func isValidSuffix(link string) bool {
	if link == "void(0)" {
		return false
	}
	if strings.HasPrefix(link, "mailto:") {
		return false
	}
	if strings.HasPrefix(link, "javascript:") {
		return false
	}

	if strings.HasPrefix(link, "#") {
		return false
	}

	return true
}

func getTitleFromPageSource(pageSource []byte) string {
	// Load the HTML document
	res := bytes.NewReader(pageSource)
	doc, err := goquery.NewDocumentFromReader(res)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	var title string
	doc.Find("title").Each(func(i int, s *goquery.Selection) {
		// 1. must have href
		// 2. concat url
		// 3. enqueue
		// 4. hub counting
		text := strings.TrimSpace(s.Text())
		re := regexp.MustCompile("(\n|\t|\r|[[:space:]][[:space:]]+)")
		text = re.ReplaceAllString(text, " ")
		title = text
		// log.Println("title", i)
		// log.Println("orig", s.Text())
		// log.Println("title trimmed", text)
	})

	if title == "" {
		doc.Find("h1").Each(func(i int, s *goquery.Selection) {
			// 1. must have href
			// 2. concat url
			// 3. enqueue
			// 4. hub counting
			text := strings.TrimSpace(s.Text())
			re := regexp.MustCompile("(\n|\t|\r|[[:space:]][[:space:]]+)")
			text = re.ReplaceAllString(text, " ")
			title = text
			// log.Println("h1", i)
			// log.Println("orig", s.Text())
			// log.Println("h1 trimmed", text)
		})
	}

	return title
}

func (manager *Manager) getNextURLList(link string, pageSource []byte) []string {
	link = strings.ToLower(strings.TrimSpace(link))
	nextURLs := []string{}

	// Load the HTML document
	res := bytes.NewReader(pageSource)
	doc, err := goquery.NewDocumentFromReader(res)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	parsedLink, err := url.Parse(link)
	if err != nil {
		log.Println("getNextURLList link", err)
		return nextURLs
	}
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		// 1. must have href
		// 2. concat url
		// 3. enqueue
		// 4. hub counting
		// 5. not self
		str, exists := s.Attr("href")
		str = strings.ToLower(strings.TrimSpace(str))
		if isValidSuffix(str) && exists {
			u, err := parsedLink.Parse(str)
			if err != nil {
				log.Println("getNextURLList href", err)
				return
			}
			// fmt.Println(i, u, strings.TrimSpace(s.Text()))
			// fmt.Println(i, u.String())
			nextURLs = append(nextURLs, u.String())
		}
	})

	// log.Println(len(nextURLs))

	return nextURLs
}
