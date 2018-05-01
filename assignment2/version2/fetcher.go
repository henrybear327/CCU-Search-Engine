package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// GetPageSource is a function that downloads page source and return it as a string
// For static pages only
func GetPageSource(url string) ([]byte, time.Duration, int) {
	// download
	startDownload := time.Now()

	res, err := http.Get(url)
	elapsedDownload := time.Since(startDownload)
	log.Printf("Downloading %s took %s", url, elapsedDownload)

	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	// get page source
	startRead := time.Now()

	robots, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	elapsedRead := time.Since(startRead)
	log.Printf("Extracting page source of %s took %s", url, elapsedRead)
	return robots, elapsedDownload + elapsedRead, res.StatusCode
}
