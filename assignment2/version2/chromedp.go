// Command simple is a chromedp example demonstrating how to do a simple google
// search.
package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/runner"
)

func getPageSource(urlstr string, title *string, pageSource *string) chromedp.Action {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.Sleep(conf.System.MinFetchTimeInterval * time.Second),
		chromedp.Title(title),
		chromedp.OuterHTML("html", pageSource),
	}
}

func getScreenshotAndPageSource(urlstr string, title *string, pageSource *string) chromedp.Action {
	var buf []byte
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.Sleep(conf.System.MinFetchTimeInterval * time.Second),
		chromedp.Title(title),
		chromedp.OuterHTML("html", pageSource),
		chromedp.CaptureScreenshot(&buf),
		chromedp.ActionFunc(func(context.Context, cdp.Executor) error {
			return ioutil.WriteFile(conf.Output.ScreenshotPath+"/"+strings.Replace(urlstr, "/", " ", -1)+".png", buf, 0644)
		}),
	}
}

func gopherGo(ctxt context.Context, pool *chromedp.Pool, urlstr string, boundedWaiting chan bool) {
	defer func(boundedWaiting chan bool) {
		<-boundedWaiting
	}(boundedWaiting)

	// allocate
	c, err := pool.Allocate(ctxt,
		runner.Flag("headless", conf.Chromedp.HeadlessMode),
		runner.Flag("no-default-browser-check", true),
		runner.Flag("no-first-run", true),
		// runner.Flag("no-sandbox", true),
		runner.ExecPath("google-chrome"))
	if err != nil {
		log.Printf("allocate url `%s` error: %v", urlstr, err)
		return
	}
	defer c.Release()

	// run tasks
	var title, pageSource string
	if conf.Output.SaveScreenshot {
		err = c.Run(ctxt, getScreenshotAndPageSource(urlstr, &title, &pageSource))
	} else {
		err = c.Run(ctxt, getPageSource(urlstr, &title, &pageSource))
	}
	fmt.Println("Back", urlstr, title)
	if err != nil {
		log.Printf("screenshot url `%s` error: %v", urlstr, err)
		// return // let the save html file continue
	}
	saveHTMLFileFromString(strings.Replace(urlstr, "/", " ", -1)+".html", pageSource)
}

// https://github.com/chromedp/examples/blob/master/pool/main.go
func getDynamicSitePageSource(link chan string, done chan bool) {
	defer func(done chan bool) {
		done <- true
	}(done)

	var err error

	// create context
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create pool
	pool, err := chromedp.NewPool()
	if err != nil {
		log.Fatalln("New pool", err)
	}
	defer pool.Shutdown()

	// loop over the URLs
	boundedWaiting := make(chan bool, conf.Chromedp.MaxConcurrentJobs)
	timeout := time.After(1 * time.Minute)
	for {
		select {
		case nextLink := <-link:
			boundedWaiting <- true
			log.Println("gopherGo", nextLink)
			go gopherGo(ctxt, pool, nextLink, boundedWaiting)
		case <-timeout:
			fmt.Println("timeout! Ending chromedp goroutine")
			return
		}
		time.Sleep(100 * time.Millisecond)
	}

	// shutdown pool
	// err = pool.Shutdown()
	// if err != nil {
	// 	log.Fatalln("pool shutdown", err)
	// }
}
