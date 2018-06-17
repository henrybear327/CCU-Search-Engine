package main

import (
	"log"
	"os"
	"strings"

	"github.com/go-ego/gse"
)

// Option is the global conf struct
type Option struct {
	segmenter Segmentation
	storage   Storage
}

func (option *Option) init() {
	option.segmenter.init()

	option.storage.init() // init map
	option.storage.load() // load data
}

/* Start Storage */
// networking-based
// load dumped data from file

type storageJSON struct {
}

func (storage *storageJSON) init() {
	storageInit()
}

func (storage *storageJSON) load() {
	// load key-value pairs from dumped file
}

// debug purpose
// load data straight from files in the folder
type storageFromFolder struct {
	folderName string
}

func (storage *storageFromFolder) init() {
	storageInit()
}

func (storage *storageFromFolder) load() {
	go func(dir string) {
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
	}(storage.folderName)
}

/* End Storage */

/* Start Segmentation */
// gse
type segmentationGSE struct {
	segmenter gse.Segmenter
}

func (s *segmentationGSE) init() {
	s.segmenter.LoadDict()
}

func (s *segmentationGSE) getSegmentedText(text []byte) []string {
	segments := s.segmenter.Segment(text)

	return gse.ToSlice(segments, false)
}

/* End Segmentation */
