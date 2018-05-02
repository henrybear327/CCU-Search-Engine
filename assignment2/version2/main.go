package main

import (
	"flag"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/BurntSushi/toml"

	_ "net/http/pprof"
)

type config struct {
	Site siteConfig // fuck it, need to start with an upper-case letter
}

type siteConfig struct {
	AlexaTopSitesURL string
}

func parseConfigFile(conf *config) {
	var configFilePath string
	flag.StringVar(&configFilePath, "configFile", "config.toml", "path to toml config file")
	flag.Parse()
	log.Println("Using config file", configFilePath)

	md, err := toml.DecodeFile(configFilePath, conf)
	if err != nil {
		log.Fatalln("Parsing config file error", err)
	}
	log.Printf("Undecoded keys: %q\n", md.Undecoded())
	log.Println("Seed", conf.Site.AlexaTopSitesURL)
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func main() {
	// profiler
	flag.Parse()
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

	// decode config file
	var conf config
	parseConfigFile(&conf)

	// scheduler starts here!
	pageSource, _, statusCode := GetStaticSitePageSource(conf.Site.AlexaTopSitesURL)
	// fmt.Println(pageSource, elapsed, statusCode)
	if statusCode == 200 {
		topURLList := ParseAlexaTopSites(pageSource)
		log.Println("Total seeding sites", len(topURLList))
	} else {
		log.Fatalln("Top Alexa sites can't be parsed!")
	}

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
