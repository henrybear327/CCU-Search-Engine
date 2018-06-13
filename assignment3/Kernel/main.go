package main

import (
	"flag"
	"log"
)

func parse() *string {
	source := flag.String("source", "docs", "the source folder to index")
	flag.Parse()

	return source
}

func check(functionName string, err error) {
	if err != nil {
		log.Fatalln("error from", functionName, err)
	}
}

func main() {
	source := parse()

	log.Println("Currently, persistent data structure is not supported")

	// build index
	if source != nil {
		indexFromDirectory(*source)
	} else {
		log.Fatalln("No source folder to load")
	}

	// run user interface
	// ui()
}
