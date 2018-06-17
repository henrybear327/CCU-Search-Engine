package main

import "log"

func check(functionName string, err error) {
	if err != nil {
		log.Fatalln("error from", functionName, err)
	}
}

var (
	option Option
	// coarse grain
	invertedIndex map[string]map[int]bool
	indexedFiles  map[int]document
)

func main() {
	/* parse command line */
	source, port := parse()

	/* config */
	option.segmenter = &segmentationGSE{}
	option.storage = &storageInitFromFolder{folderName: source}

	/* init */
	option.init()

	// run command line interface for searching
	go searchCLI()

	// let's start!
	listen(port) // blocking!
}
