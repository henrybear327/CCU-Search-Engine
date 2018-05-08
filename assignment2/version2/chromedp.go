// Command simple is a chromedp example demonstrating how to do a simple google
// search.
package main

import (
	"context"
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
		chromedp.Sleep(conf.System.minFetchTimeDuration),
		chromedp.Title(title),
		chromedp.OuterHTML("html", pageSource),
	}
}

func getScreenshotAndPageSource(urlstr string, title *string, pageSource *string /*, nodes *[]*cdp.Node*/) chromedp.Action {
	var buf []byte
	createFolderIfNotExist(conf.Output.ScreenshotPath)
	path := conf.Output.ScreenshotPath + "/" + getTopLevelDomain(urlstr)
	createFolderIfNotExist(path)

	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.Sleep(conf.System.minFetchTimeDuration),
		chromedp.Title(title),
		chromedp.OuterHTML("html", pageSource),
		chromedp.CaptureScreenshot(&buf),
		chromedp.ActionFunc(func(context.Context, cdp.Executor) error {
			return ioutil.WriteFile(path+"/"+strings.Replace(urlstr, "/", " ", -1)+".png", buf, 0644)
		}),
		// chromedp.Nodes(`a`, nodes, chromedp.ByQueryAll),
	}
}

type dynamicFetchingDataQuery struct {
	link          string
	resultChannel chan dynamicFetchingDataResult
}

type dynamicFetchingDataResult struct {
	title           string
	pageSource      string
	requiresRestart bool
}

func gopherGo(ctxt context.Context, pool *chromedp.Pool, query dynamicFetchingDataQuery, semaphore chan bool) {
	defer func(semaphore chan bool) {
		<-semaphore
	}(semaphore)

	result := dynamicFetchingDataResult{requiresRestart: false}

	// allocate
	c, err := pool.Allocate(ctxt,
		runner.Flag("headless", conf.Chromedp.HeadlessMode),
		runner.Flag("no-default-browser-check", true),
		runner.Flag("no-first-run", true),
		// runner.Flag("no-sandbox", true),
		runner.ExecPath(conf.Chromedp.ExecPath),
	)
	if err != nil {
		log.Printf("[error] allocate url `%s` error: %v", query.link, err)

		result.requiresRestart = true
		query.resultChannel <- result
		return
	}
	defer c.Release()

	// run tasks
	// var nodes []*cdp.Node
	var title, pageSource string
	if conf.Output.SaveScreenshot {
		err = c.Run(ctxt, getScreenshotAndPageSource(query.link, &title, &pageSource /*, &nodes*/))
	} else {
		err = c.Run(ctxt, getPageSource(query.link, &title, &pageSource))
	}
	if err != nil {
		log.Printf("[error] chromedp back `%s` error: %v", query.link, err)
		// return // let the save html file continue
	}
	result.title = title
	result.pageSource = pageSource

	log.Println("Back", query.link, title)

	query.resultChannel <- result

	// fmt.Println(len(nodes))
	// for _, n := range nodes {
	// 	fmt.Println(n.AttributeValue("href"))
	// }
}

// https://github.com/chromedp/examples/blob/master/pool/main.go
func getDynamicSitePageSource(data chan dynamicFetchingDataQuery) {
	var err error

	// create context
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create pool
	fileName := "chromedp.log"
	logFile, err := os.Create(fileName)
	defer logFile.Close()
	if err != nil {
		log.Fatalln("[error] chromedp log", err)
	}

	debugLog := log.New(logFile, "[Chromedp]", log.LstdFlags)
	debugLog.SetFlags(debugLog.Flags() | log.LstdFlags)
	pool, err := chromedp.NewPool(chromedp.PoolLog(debugLog.Printf, debugLog.Printf, debugLog.Printf))
	if err != nil {
		log.Fatalln("[error] New pool", err)
	}
	defer func() {
		err := pool.Shutdown()
		if err != nil {
			log.Println("[error] defer shutdown", err)
		}
	}()

	// loop over the URLs
	semaphore := make(chan bool, conf.Chromedp.MaxConcurrentJobs)
	timeout := time.After(conf.System.maxRunningTimeDuration)
	for {
		select {
		case nextData := <-data:
			semaphore <- true
			log.Println("gopherGo", nextData.link)
			go gopherGo(ctxt, pool, nextData, semaphore)
		case <-timeout:
			log.Println("[error] Chromedp timeout! Ending chromedp goroutine")
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// func testSelectAll() {
// 	// create context
// 	ctxt, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	c, err := chromedp.New(ctxt, chromedp.WithRunnerOptions(
// 		runner.Flag("headless", true),
// 		runner.Flag("disable-gpu", true)))

// 	var nodes []*cdp.Node
// 	t := chromedp.Tasks{
// 		chromedp.Navigate("https://www.npr.org/"),
// 		chromedp.Sleep(time.Second * 2),
// 		chromedp.Nodes("a", &nodes, chromedp.ByQueryAll),
// 	}
// 	err = c.Run(ctxt, t)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println(len(nodes))
// 	for _, n := range nodes {
// 		fmt.Println(n.AttributeValue("href"))
// 	}
// }
