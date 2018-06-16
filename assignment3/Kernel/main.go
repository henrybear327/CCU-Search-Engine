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

	// build index
	indexFromDirectory(*source)

	// run user interface
	// ui()
}
