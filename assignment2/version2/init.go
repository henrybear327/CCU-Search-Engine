package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"net/url"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	"golang.org/x/net/publicsuffix"
)

type config struct {
	Site     siteConfig // fuck it, need to start with an upper-case letter
	Output   outputConfig
	System   systemConfig
	Chromedp chromedpConfig
	MongoDB  mongoDBConfig
}

type siteConfig struct {
	UseAlexaTopSites bool
	AlexaTopSitesURL string
	ManualSeedFile   string
}

type outputConfig struct {
	Seedfile          string
	ParsingResultFile string
	ScreenshotPath    string
	PageSourcePath    string

	SlowAction         string
	SlowActionDuration time.Duration

	SaveScreenshot bool
}

type systemConfig struct {
	MaxDistinctPagesToFetchPerSite int
	MinFetchTimeInterval           string
	minFetchTimeDuration           time.Duration
	MaxRunningTime                 string
	maxRunningTimeDuration         time.Duration
	MaxGoRountinesPerSite          int
}

type chromedpConfig struct {
	HeadlessMode      bool
	MaxConcurrentJobs int
	ExecPath          string
}

type mongoDBConfig struct {
	URL      string
	Database string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func startCPUProfiling(cpuprofile *string) {
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}
}

func startMemProfiling(memprofile *string) {
	// profiler
	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
		f.Close()
	}
}

func parseConfigFile() {
	var configFilePath string
	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
	var memprofile = flag.String("memprofile", "", "write memory profile to `file`")
	flag.StringVar(&configFilePath, "configFile", "config.toml", "path to toml config file")
	flag.Parse()

	log.Println("Using config file", configFilePath)
	startCPUProfiling(cpuprofile)
	startMemProfiling(memprofile)

	md, err := toml.DecodeFile(configFilePath, &conf)
	if err != nil {
		log.Fatalln("Parsing config file error", err)
	}
	log.Printf("Undecoded keys: %q\n", md.Undecoded())

	log.Println("=====config file=====")
	log.Println("UseAlexaTopSites", conf.Site.UseAlexaTopSites)
	log.Println("AlexaTopSitesURL", conf.Site.AlexaTopSitesURL)
	log.Println("ManualSeedURL", conf.Site.ManualSeedFile)
	log.Println("MaxDistinctPagesToFetchPerSite", conf.System.MaxDistinctPagesToFetchPerSite)
	log.Println("MinFetchTimeInterval", conf.System.MinFetchTimeInterval)
	log.Println("MaxRunningTime", conf.System.MaxRunningTime)
	log.Println("MaxGoRountinesPerSite", conf.System.MaxGoRountinesPerSite)
	log.Println("Seedfile", conf.Output.Seedfile)
	log.Println("ParsingResultFile", conf.Output.ParsingResultFile)
	log.Println("ScreenshotPath", conf.Output.ScreenshotPath)
	log.Println("PageSourcePath", conf.Output.PageSourcePath)
	log.Println("SlowAction", conf.Output.SlowAction)
	log.Println("SaveScreenshot", conf.Output.SaveScreenshot)
	log.Println("HeadlessMode", conf.Chromedp.HeadlessMode)
	log.Println("MaxConcurrentJobs", conf.Chromedp.MaxConcurrentJobs)
	log.Println("ExecPath", conf.Chromedp.ExecPath)
	log.Println("MongoDB Database", conf.MongoDB.Database)
	log.Println("MongoDB URL", conf.MongoDB.URL)
	log.Println("=====================")
	{
		var err error
		conf.System.minFetchTimeDuration, err = time.ParseDuration(conf.System.MinFetchTimeInterval)
		if err != nil {
			log.Fatalln("MinFetchTimeInterval", err)
		}

		conf.System.maxRunningTimeDuration, err = time.ParseDuration(conf.System.MaxRunningTime)
		if err != nil {
			log.Fatalln("MaxRunningTime", err)
		}

		conf.Output.SlowActionDuration, err = time.ParseDuration(conf.Output.SlowAction)
		if err != nil {
			log.Fatalln("SlowActionDuration", err)
		}
	}
}

func getSeedSites() ([]string, []string) {
	var seedSiteList []string
	var seedSiteOptionList []string

	if conf.Site.UseAlexaTopSites {
		pageSource, statusCode := getStaticSitePageSource(conf.Site.AlexaTopSitesURL)
		// log.Println(pageSource, statusCode)
		if statusCode == 200 {
			seedSiteList = parseAlexaTopSites(pageSource)
			log.Println("Total seeding sites", len(seedSiteList))
		} else {
			log.Fatalln("Top Alexa sites can't be parsed!")
		}
	} else {
		f, err := os.Open(conf.Site.ManualSeedFile)
		check(err)
		// seedSiteList = append(seedSiteList, conf.Site.ManualSeedURL...) // cool

		r := bufio.NewReader(f)

		for {
			str, _, err := r.ReadLine()
			if err == io.EOF {
				break
			}
			check(err)

			line := string(str)
			lineComponents := strings.Split(line, " ")
			if len(lineComponents) != 2 {
				log.Fatalln(line, "in the seed site file is not accepted")
			}

			seedSiteList = append(seedSiteList, lineComponents[0])
			seedSiteOptionList = append(seedSiteOptionList, lineComponents[1])
		}
	}

	outputSeedingSites(seedSiteList, seedSiteOptionList)
	return seedSiteList, seedSiteOptionList
}

func getTopLevelDomain(link string) string {
	link = strings.ToLower(strings.TrimSpace(link))

	u, err := url.Parse(link)
	if err != nil {
		log.Fatalln("Parsing hostname error")
	}
	host := u.Hostname()

	linkTLD, err := publicsuffix.EffectiveTLDPlusOne(host)
	if err != nil {
		log.Println("isExternalSite EffectiveTLDPlusOne err", err)
		return ""
	}
	return linkTLD
}

func prepareSeedSites(seedSiteList []string, seedSiteOption []string) map[string]*Manager {
	startDownload := time.Now()

	totalSites := len(seedSiteList)
	totalOptions := len(seedSiteOption)
	managers := make(map[string]*Manager) // pointer bug!!
	done := make(chan bool, totalSites)

	for i, link := range seedSiteList {
		log.Println("Starting", link, "preprocessing", "TLD =", getTopLevelDomain(link))

		useStaticLoad := true
		if totalOptions == totalSites && seedSiteOption[i] == "@useChrome" {
			useStaticLoad = false
		}
		log.Println(link, "useStaticLoad", useStaticLoad)

		newManager := Manager{
			link:            link,
			urlQueueLock:    new(sync.RWMutex),
			urlInQueueLock:  new(sync.RWMutex),
			urlFetchedLock:  new(sync.RWMutex),
			urlFetched:      make(map[string]bool),
			urlInQueue:      make(map[string]bool),
			tld:             getTopLevelDomain(link),
			useLinksFromXML: false,
			useStaticLoad:   useStaticLoad}

		managers[link] = &newManager

		cur := managers[link]
		// log.Println("tld", link, cur.tld)
		go cur.preprocess(done)
	}

	for i := 0; i < totalSites; i++ {
		<-done
	}

	// check for sites that are simply not allowed to parse
	log.Println("Manager count before removing", len(managers))
	for _, link := range seedSiteList {
		cur := managers[link]

		if cur.isBannedByRobotTXT(link) {
			log.Println(link, "is banned from being parsed")
			delete(managers, link)
		} else {
			if cur.robot != nil {
				cur.crawlDelay = cur.robot.FindGroup("CCU-Assignment-Bot").CrawlDelay
			} else {
				cur.crawlDelay = 0
			}
		}
	}
	log.Println("Manager count after removing", len(managers))

	elapsedDownload := time.Since(startDownload)
	log.Println("Total preprocessing time", elapsedDownload)
	return managers
}
