package main

import (
	"flag"
	"log"

	"github.com/BurntSushi/toml"
)

type config struct {
	AlexaTopSitesURL string
}

func parseConfigFile(conf *config) {
	var configFilePath string
	flag.StringVar(&configFilePath, "configFile", "config.toml", "path to toml config file")
	log.Println("Using config file", configFilePath)
	flag.Parse()

	if _, err := toml.DecodeFile(configFilePath, conf); err != nil {
		log.Fatalln("Parsing config file error", err)
	}
	log.Println("Seed", conf.AlexaTopSitesURL)
}

func main() {
	// decode config file
	var conf config
	parseConfigFile(&conf)

	// scheduler starts here!
	pageSource, _, statusCode := GetStaticSitePageSource(conf.AlexaTopSitesURL)
	// fmt.Println(pageSource, elapsed, statusCode)
	if statusCode == 200 {
		topURLList := ParseAlexaTopSites(pageSource)
		log.Println("Total seeding sites", len(topURLList))
	} else {
		log.Fatalln("Top Alexa sites can't be parsed!")
	}
}
