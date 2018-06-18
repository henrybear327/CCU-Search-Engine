package main

import "log"

func check(functionName string, err error) {
	if err != nil {
		log.Fatalln("error from", functionName, err)
	}
}

var (
	configuration Configuration

	nextDocID int

	// coarse grain
	invertedIndex invertedIndexData
	indexedFiles  indexedFilesData
)

func main() {
	/* parse command line */
	source, port := parse()

	/* config */
	configuration.segmenter = &segmentationGSE{}
	configuration.storage = &storageFromFolder{folderName: source}

	/* init */
	nextDocID = 0
	configuration.init()

	// run command line interface for searching
	go searchCLI()

	// let's start!
	listen(port) // blocking!
}
