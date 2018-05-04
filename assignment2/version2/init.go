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

	"github.com/BurntSushi/toml"
)

type config struct {
	Site   siteConfig // fuck it, need to start with an upper-case letter
	Output outputConfig
	System systemConfig
}

type siteConfig struct {
	UseAlexaTopSites bool
	AlexaTopSitesURL string
	ManualSeedFile   string
}

type outputConfig struct {
	Seedfile          string
	ParsingResultFile string

	SlowAction int64
}

type systemConfig struct {
	MaxDistinctPagesToFetchPerSite int
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
	log.Println("Seedfile", conf.Output.Seedfile)
	log.Println("ParsingResultFile", conf.Output.ParsingResultFile)
	log.Println("SlowAction", conf.Output.SlowAction)
	log.Println("=====================")
}

func getSeedSites() []string {
	var seedSiteList []string

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
			seedSiteList = append(seedSiteList, string(str))
		}
	}

	outputSeedingSites(seedSiteList)
	return seedSiteList
}

func prepareSeedSites(seedSiteList []string) map[string]*Manager {
	totalSites := len(seedSiteList)
	managers := make(map[string]*Manager) // pointer bug!!
	done := make(chan bool, totalSites)

	for _, link := range seedSiteList {
		u, err := url.Parse(link)
		if err != nil {
			log.Fatalln("Parsing hostname error")
		}
		host := u.Hostname()

		newManager := Manager{
			link:            link,
			urlQueueLock:    new(sync.RWMutex),
			urlInQueueLock:  new(sync.RWMutex),
			urlFetchedLock:  new(sync.RWMutex),
			urlFetched:      make(map[string]bool),
			urlInQueue:      make(map[string]bool),
			tld:             strings.ToLower(getTLD(host)),
			useLinksFromXML: false}

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
		// fmt.Println(cur.link)
		if cur.isBannedByRobotTXT(link) {
			log.Println(link, "is banned from being parsed")
			delete(managers, link)
		}
	}
	log.Println("Manager count after removing", len(managers))
	return managers
}
