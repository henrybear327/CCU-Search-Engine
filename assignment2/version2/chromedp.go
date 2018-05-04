// Command simple is a chromedp example demonstrating how to do a simple google
// search.
package main

import (
	"context"
	"io/ioutil"
	"log"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/client"
)

func getDynamicSitePageSource(link string, done chan bool) {
	log.Println("Dynamic fetch", link)
	var err error

	// create context (place to fill data)
	ctxt, cancel := context.WithCancel(context.Background())
	// ctxt, cancel := context.WithTimeout(context.Background(), (conf.System.MinFetchTimeInterval+5)*time.Second) // doesn't work for cnn.com GG
	defer cancel()
	defer func(done chan bool) {
		done <- true
	}(done)

	var c *chromedp.CDP
	if conf.Chromedp.HeadlessMode {
		var err error
		// c, err = chromedp.New(ctxt, chromedp.WithTargets(client.New().WatchPageTargets(ctxt)), chromedp.WithLog(log.Printf))
		c, err = chromedp.New(ctxt, chromedp.WithTargets(client.New().WatchPageTargets(ctxt)))
		if err != nil {
			log.Println("headless mode", err)
		}
	}

	if conf.Chromedp.HeadlessMode == false || c == nil {
		var err error
		// create chrome instance (new browser)
		// c, err = chromedp.New(ctxt, chromedp.WithLog(log.Printf))
		c, err = chromedp.New(ctxt)
		if err != nil {
			log.Fatalln("non headless mode", err)
		}
	}

	// run task list
	var title, pageSource string
	err = c.Run(ctxt, getTitleAndPageSourceFromLink(link, &title, &pageSource))
	if err != nil {
		log.Println("getTitleAndPageSourceFromLink", err)
	}
	// fmt.Println(title)
	// fmt.Println(pageSource)
	saveHTMLFileFromString(title+".html", pageSource)

	if conf.Chromedp.HeadlessMode == false {
		// shutdown chrome
		err = c.Shutdown(ctxt)
		if err != nil {
			log.Fatalln("chromedp Shutdown", err)
		}

		// wait for chrome to finish
		err = c.Wait()
		if err != nil {
			log.Fatalln("chromedp wait", err)
		}
	}
}

func getTitleAndPageSourceFromLink(link string, title *string, pageSource *string) chromedp.Tasks {
	var buf []byte
	return chromedp.Tasks{
		chromedp.Navigate(link),
		// chromedp.WaitVisible(`#hplogo`, chromedp.ByID),
		chromedp.Sleep(conf.System.MinFetchTimeInterval * time.Second),
		chromedp.Title(title),
		chromedp.CaptureScreenshot(&buf),
		chromedp.ActionFunc(func(context.Context, cdp.Executor) error {
			return ioutil.WriteFile(*title+".png", buf, 0644)
		}),
		chromedp.OuterHTML("html", pageSource),
	}
}
