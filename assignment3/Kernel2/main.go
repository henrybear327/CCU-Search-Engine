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
	debugPtr := flag.Bool("debug", false, "debug mode")
	flag.Parse()

	debug = *debugPtr
	return *port
}

func main() {
	// setup
	port := parseCommandLine()
	// log.Println("port", port, "debug", debug)

	seg.init()
	idxer.init()

	// start
	listen(port)
}
