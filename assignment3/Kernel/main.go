package main

var (
	config configuration
)

func main() {
	/* parse command line */
	source, port, gobFile, debugMode := parseCommandLine()

	/* config */
	// setup interfaces
	config.segmenter = &segmentationGSE{}

	if debugMode == 0 {
		config.storage = &storageStupid{folderName: source}
	} else if debugMode == 1 {
		config.storage = &storageV1{}
	}

	/* init */
	config.init(gobFile)

	/* run command line interface for searching in the background */
	go searchCLI()

	/* TODO: handle ctrl+c interrupt for data saving */

	/* let's start! */
	listen(port) // blocking!
}
