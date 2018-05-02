package main

import (
	"fmt"
	_ "net/http/pprof"
)

func main() {
	var conf config
	parseConfigFile(&conf)

	// scheduler starts here!
	seedSiteList := getSeedSites(&conf)
	managers := prepareSeedSites(seedSiteList)

	fmt.Println("Manager count", len(managers))
}
