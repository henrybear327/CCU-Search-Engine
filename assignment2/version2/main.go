package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/BurntSushi/toml"
	"github.com/temoto/robotstxt"

	_ "net/http/pprof"
)

type config struct {
	Site   siteConfig // fuck it, need to start with an upper-case letter
	Output outputConfig
}

type siteConfig struct {
	UseAlexaTopSites bool
	AlexaTopSitesURL string
	ManualSeedFile   string
}

type outputConfig struct {
	Seedfile string
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

func parseConfigFile(conf *config) {
	var configFilePath string
	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
	var memprofile = flag.String("memprofile", "", "write memory profile to `file`")
	flag.StringVar(&configFilePath, "configFile", "config.toml", "path to toml config file")
	flag.Parse()

	log.Println("Using config file", configFilePath)
	startCPUProfiling(cpuprofile)
	startMemProfiling(memprofile)

	md, err := toml.DecodeFile(configFilePath, conf)
	if err != nil {
		log.Fatalln("Parsing config file error", err)
	}
	log.Printf("Undecoded keys: %q\n", md.Undecoded())

	log.Println("=====config file=====")
	log.Println("UseAlexaTopSites", conf.Site.UseAlexaTopSites)
	log.Println("AlexaTopSitesURL", conf.Site.AlexaTopSitesURL)
	log.Println("ManualSeedURL", conf.Site.ManualSeedFile)
	log.Println("=====================")
}

func getSeedSites(conf *config) []string {
	var seedSiteList []string

	if conf.Site.UseAlexaTopSites {
		pageSource, statusCode := GetStaticSitePageSource(conf.Site.AlexaTopSitesURL)
		// fmt.Println(pageSource, statusCode)
		if statusCode == 200 {
			seedSiteList = ParseAlexaTopSites(pageSource)
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

	OutputSeedingSites(seedSiteList, conf)
	return seedSiteList
}

// RobotData for channel communication
type RobotData struct {
	url   string
	robot *robotstxt.RobotsData
}

func prepareSeedSites(seedSiteList []string) {
	totalSites := len(seedSiteList)
	done := make(chan RobotData, totalSites)

	for _, url := range seedSiteList {
		// parse robots.txt
		go ParseRobotsTxt(url, done)

		// parse sitemap.xml

		// prepare queue

	}

	robotsCollection := make(map[string]*robotstxt.RobotsData)
	// cnt := 0
	for i := 0; i < totalSites; i++ {
		ret := <-done
		robotsCollection[ret.url] = ret.robot

		// if ret.robot != nil && len(ret.robot.Sitemaps) > 0 {
		// 	cnt++
		// }
	}
	// fmt.Println(len(robotsCollection))
	// fmt.Println("sitemap", cnt)
}

func main() {
	// decode config file
	var conf config
	parseConfigFile(&conf)

	// scheduler starts here!
	seedSiteList := getSeedSites(&conf)
	prepareSeedSites(seedSiteList)
}
