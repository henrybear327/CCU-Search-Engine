package main

var (
	config configuration

	nextDocID int

	invertedIndex invertedIndexData
	indexedFiles  indexedFilesData
)

func main() {
	/* parse command line */
	source, port := parseCommandLine()

	/* config */
	// setup interfaces
	config.segmenter = &segmentationGSE{}
	config.storage = &storageFromFolder{folderName: source}

	/* init */
	nextDocID = 0
	config.init()

	/* run command line interface for searching in the background */
	go searchCLI()

	/* TODO: handle ctrl+c interrupt for data saving */

	/* let's start! */
	listen(port) // blocking!
}
