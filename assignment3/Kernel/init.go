package main

import (
	"flag"
)

func parse() (string, int) {
	source := flag.String("source", "docs", "the source folder to index")
	port := flag.Int("port", 8001, "port to listen for requests")
	flag.Parse()

	return *source, *port
}

func loadDataFromFile() {

}
