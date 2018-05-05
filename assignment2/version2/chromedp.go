// Command simple is a chromedp example demonstrating how to do a simple google
// search.
package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/runner"
)

func getPageSource(urlstr string, title *string, pageSource *string) chromedp.Action {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.Sleep(conf.System.MinFetchTimeInterval),
		chromedp.Title(title),
		chromedp.OuterHTML("html", pageSource),
	}
}

func getScreenshotAndPageSource(urlstr string, title *string, pageSource *string, nodes *[]*cdp.Node) chromedp.Action {
	var buf []byte
	ret := chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.Sleep(conf.System.MinFetchTimeInterval),
		chromedp.Title(title),
		chromedp.OuterHTML("html", pageSource),
		chromedp.CaptureScreenshot(&buf),
		chromedp.ActionFunc(func(context.Context, cdp.Executor) error {
			return ioutil.WriteFile(conf.Output.ScreenshotPath+"/"+strings.Replace(urlstr, "/", " ", -1)+".png", buf, 0644)
		}),
		chromedp.Nodes(`a`, nodes, chromedp.ByQueryAll),
	}

	return ret
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
	var nodes []*cdp.Node
	if conf.Output.SaveScreenshot {
		err = c.Run(ctxt, getScreenshotAndPageSource(urlstr, &title, &pageSource, &nodes))
	} else {
		err = c.Run(ctxt, getPageSource(urlstr, &title, &pageSource))
	}
	fmt.Println("Back", urlstr, title)
	if err != nil {
		log.Printf("screenshot url `%s` error: %v", urlstr, err)
		// return // let the save html file continue
	}
	saveHTMLFileFromString(strings.Replace(urlstr, "/", " ", -1)+".html", pageSource)

	fmt.Println(len(nodes))
	for _, n := range nodes {
		fmt.Println(n.AttributeValue("href"))
	}
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

	// c, err := chromedp.New(ctxt, chromedp.WithRunnerOptions(
	// 	runner.Flag("headless", true),
	// 	runner.Flag("disable-gpu", true)))

	// var nodes []*cdp.Node
	// t := chromedp.Tasks{
	// 	chromedp.Navigate("https://www.npr.org/"),
	// 	chromedp.Sleep(time.Second * 2),
	// 	chromedp.Nodes("a", &nodes, chromedp.ByQueryAll),
	// }
	// err = c.Run(ctxt, t)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(len(nodes))
	// for _, n := range nodes {
	// 	fmt.Println(n.AttributeValue("href"))
	// }

	// create pool
	fileName := "chromedp.log"
	logFile, err := os.Create(fileName)
	defer logFile.Close()
	if err != nil {
		log.Fatalln("chromedp log", err)
	}

	debugLog := log.New(logFile, "[Chromedp]", log.LstdFlags)
	debugLog.SetFlags(debugLog.Flags() | log.LstdFlags)
	pool, err := chromedp.NewPool(chromedp.PoolLog(debugLog.Printf, debugLog.Printf, debugLog.Printf))
	if err != nil {
		log.Fatalln("New pool", err)
	}
	defer pool.Shutdown()

	// loop over the URLs
	boundedWaiting := make(chan bool, conf.Chromedp.MaxConcurrentJobs)
	timeout := time.After(conf.System.MaxRunningTime)
	for {
		select {
		case nextLink := <-link:
			boundedWaiting <- true
			log.Println("gopherGo", nextLink)
			go gopherGo(ctxt, pool, nextLink, boundedWaiting)
		case <-timeout:
			fmt.Println("Chromedp timeout! Ending chromedp goroutine")
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
}
