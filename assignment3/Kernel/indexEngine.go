package main

import (
	"log"
	"os"
	"strings"
)

type document struct {
	filename string
}

func indexFromDirectory(dir string) {
	log.Println("Indexing directory", dir)

	directory, err := os.Open(dir) // open directory
	check("os.Open", err)
	defer directory.Close()

	filesInDirectory, err := directory.Readdir(-1) // read directory
	check("directory.Readdir", err)

	if len(filesInDirectory) == 0 { // empty
		// return fmt.Errorf("There are no files in %s", dir)
		log.Fatalf("There are no files in %s", dir)
	}

	index = make(map[string]map[int]bool)
	indexedFiles = make(map[int]document)
	for docID, file := range filesInDirectory {
		if file.IsDir() == false { // non-recursive
			if strings.HasPrefix(file.Name(), ".") {
				log.Println("skipping", file.Name())
				continue
			}

			filename := dir + "/" + file.Name()
			log.Println("indexing", filename)
			parseDocument(filename, docID)
		} else {
			log.Println("Recursive is not supported")
		}
	}

	log.Println("Indexing directory, done")
}
