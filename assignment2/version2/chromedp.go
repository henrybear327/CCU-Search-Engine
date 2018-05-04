// Command simple is a chromedp example demonstrating how to do a simple google
// search.
package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/client"
)

func run() {
	var err error

	// create context (place to fill data)
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	var c *chromedp.CDP
	if conf.Chromedp.HeadlessMode {
		var err error
		// create chrome tab
		c, err = chromedp.New(ctxt, chromedp.WithTargets(client.New().WatchPageTargets(ctxt)), chromedp.WithLog(log.Printf))
		if err != nil {
			log.Fatal(err)
		}
	} else {
		var err error
		// create chrome instance (new browser)
		c, err = chromedp.New(ctxt, chromedp.WithLog(log.Printf))
		if err != nil {
			log.Fatal(err)
		}
	}

	// run task list
	var title string
	var pageSource string
	link := "https://edition.cnn.com/2018/05/04/europe/nobel-prize-for-literature-swedish-academy-postponed-intl/index.html"
	err = c.Run(ctxt, googleSearch(link, &title, &pageSource))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(title)
	// fmt.Println(pageSource)
	saveHTMLFileFromString("sample.html", pageSource)

	if conf.Chromedp.HeadlessMode == false {
		// shutdown chrome
		err = c.Shutdown(ctxt)
		if err != nil {
			log.Fatal(err)
		}

		// wait for chrome to finish
		err = c.Wait()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func googleSearch(link string, title *string, pageSource *string) chromedp.Tasks {
	var buf []byte
	return chromedp.Tasks{
		chromedp.Navigate(link),
		// chromedp.WaitVisible(`#hplogo`, chromedp.ByID),
		chromedp.Sleep(5 * time.Second),
		chromedp.CaptureScreenshot(&buf),
		chromedp.ActionFunc(func(context.Context, cdp.Executor) error {
			return ioutil.WriteFile("screenshot.png", buf, 0644)
		}),
		chromedp.Title(title),
		chromedp.OuterHTML("html", pageSource),
	}
}
