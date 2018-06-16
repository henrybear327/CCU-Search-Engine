package main

import "log"

func check(functionName string, err error) {
	if err != nil {
		log.Fatalln("error from", functionName, err)
	}
}

func main() {
	// init
	source, port := parse()
	loadDataFromFile()

	// build index
	go indexFromDirectory(source)

	// run user interface
	go ui()

	// let's start!
	listenOnPort(port)
}
