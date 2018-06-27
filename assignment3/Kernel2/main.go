package main

import (
	"flag"
)

var (
	seg   segmenter
	idxer indexer
	debug bool
)

func parseCommandLine() int {
	port := flag.Int("port", 8001, "port to listen for requests")
	debug = *(flag.Bool("debug", false, "debug mode"))
	flag.Parse()

	return *port
}

func main() {
	// setup
	port := parseCommandLine()

	seg.init()
	idxer.init()

	// start
	listen(port)
}
