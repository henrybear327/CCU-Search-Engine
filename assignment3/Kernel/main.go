package main

import "log"

func check(functionName string, err error) {
	if err != nil {
		log.Fatalln("error from", functionName, err)
	}
}

var (
	config configuration

	nextDocID int

	invertedIndex invertedIndexData
	indexedFiles  indexedFilesData
)

func main() {
	/* parse command line */
	source, port := parse()

	/* config */
	config.segmenter = &segmentationGSE{}
	config.storage = &storageFromFolder{folderName: source}

	/* init */
	nextDocID = 0
	config.init()

	// run command line interface for searching
	go searchCLI()

	// TODO: handle ctrl+c interrupt for data saving

	// let's start!
	listen(port) // blocking!
}
