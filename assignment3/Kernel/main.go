package main

import "log"

func check(functionName string, err error) {
	if err != nil {
		log.Fatalln("error from", functionName, err)
	}
}

var (
	option Option
)

func main() {
	// config
	// setup the segmentation program to use
	option.segmenter = &segmentationGSE{}
	option.segmenter.init()

	// init
	source, port := parse()
	loadDataFromFile()

	// build index (debug only)
	go indexFromDirectory(source)

	// run user interface for searching
	go ui()

	// let's start!
	listen(port)
}
