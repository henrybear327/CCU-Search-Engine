package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

var fetchingCounter struct {
	sync.Mutex
	n int
}

// GetStaticSitePageSource is a function that downloads page source of assigned link
// and return it as a []byte
func getStaticSitePageSource(link string) ([]byte, int) {
	log.Println("getStaticSitePageSource start", link)

	fetchingCounter.Lock()
	for fetchingCounter.n >= conf.System.MaxConcurrentFetch {
		time.Sleep(100 * time.Millisecond)
		log.Println("Waiting for fetchingCounter lock")
	}
	fetchingCounter.n++
	fetchingCounter.Unlock()
	defer func() {
		fetchingCounter.Lock()
		fetchingCounter.n--
		fetchingCounter.Unlock()
	}()

	link = strings.TrimSpace(link)
	// download
	startDownload := time.Now()

	res, err := http.Get(link)
	elapsedDownload := time.Since(startDownload)
	if elapsedDownload >= conf.Output.SlowActionDuration {
		log.Printf("Downloading %s took %s", link, elapsedDownload)
	}

	if err != nil {
		color.Set(color.FgRed)
		log.Println("[error] GetStaticSitePageSource http.Get", err)
		color.Unset()
		return make([]byte, 0), -1
	}
	defer res.Body.Close()

	// get page source
	startRead := time.Now()

	robots, err := ioutil.ReadAll(res.Body)
	if err != nil {
		color.Set(color.FgRed)
		log.Println("[error] GetStaticSitePageSource ioutil.ReadAll", err)
		color.Unset()
		return make([]byte, 0), -1
	}

	elapsedRead := time.Since(startRead)
	if elapsedRead >= conf.Output.SlowActionDuration {
		log.Printf("Extracting page source of %s took %s", link, elapsedRead)
	}

	log.Println("getStaticSitePageSource end", link, res.StatusCode)
	return robots, res.StatusCode
}
