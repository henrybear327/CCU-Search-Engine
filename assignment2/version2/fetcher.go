package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// GetStaticSitePageSource is a function that downloads page source of assigned url
// and return it as a []byte
func GetStaticSitePageSource(url string) ([]byte, int) {
	// download
	startDownload := time.Now()

	res, err := http.Get(url)
	elapsedDownload := time.Since(startDownload)
	log.Printf("Downloading %s took %s", url, elapsedDownload)

	if err != nil {
		log.Fatalln("GetStaticSitePageSource http.Get", err)
	}
	defer res.Body.Close()

	// get page source
	startRead := time.Now()

	robots, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln("GetStaticSitePageSource ioutil.ReadAll", err)
	}

	elapsedRead := time.Since(startRead)
	log.Printf("Extracting page source of %s took %s", url, elapsedRead)
	return robots, res.StatusCode
}
