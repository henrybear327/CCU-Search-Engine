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
	pool, err := chromedp.NewPool( /*chromedp.PoolLog(log.Printf, log.Printf, log.Printf)*/ )
	if err != nil {
		log.Fatalln("New pool", err)
	}
	defer pool.Shutdown()

	// loop over the URLs
	boundedWaiting := make(chan bool, 15)
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
	// var buf []byte
	var title, pageSource string
	err = c.Run(ctxt, getScreenshotAndPageSource(urlstr /*&buf,*/, &title, &pageSource))
	fmt.Println("Back", urlstr, title)
	if err != nil {
		log.Printf("screenshot url `%s` error: %v", urlstr, err)
		// return // TODO ??
	}
	saveHTMLFileFromString(strings.Replace(urlstr, "/", " ", -1)+".html", pageSource)
}

func getScreenshotAndPageSource(urlstr string /*, picbuf *[]byte*/, title *string, pageSource *string) chromedp.Action {
	var buf []byte
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.Sleep(conf.System.MinFetchTimeInterval * time.Second),
		chromedp.Title(title),
		// chromedp.ActionFunc(func(ctxt context.Context, h cdp.Executor) error {
		// 	buf, err := page.CaptureScreenshot().Do(ctxt, h)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	*picbuf = buf
		// 	return ioutil.WriteFile(*title+".png", buf, 0644)
		// }),
		chromedp.CaptureScreenshot(&buf),
		chromedp.ActionFunc(func(context.Context, cdp.Executor) error {
			return ioutil.WriteFile(conf.Output.ScreenshotPath+"/"+strings.Replace(urlstr, "/", " ", -1)+".png", buf, 0644)
		}),
		chromedp.OuterHTML("html", pageSource),
	}
}
