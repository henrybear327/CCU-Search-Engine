package main

import (
	"fmt"
	"log"
)

func main() {
	pageSource, _, statusCode := GetStaticSitePageSource("https://www.alexa.com/topsites/countries/TW")
	// fmt.Println(pageSource, elapsed, statusCode)
	if statusCode == 200 {
		topURLList := ParseAlexaTopSites(pageSource)
		fmt.Println(len(topURLList))
	} else {
		log.Fatalln("Top Alexa sites can't be parsed!")
	}
}
