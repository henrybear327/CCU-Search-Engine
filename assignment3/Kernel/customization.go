package main

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/go-ego/gse"
)

/* Start Storage */

/* File version */
type storageStupid struct {
	folderName string
}

func parseFromFile(filename string, docID int) map[string][]int {
	log.Println("indexing", filename)

	// load document
	// open file
	f, err := os.Open(filename)
	check("os.Open", err)
	defer f.Close()

	// scan lines, one by one
	r := bufio.NewReader(f)

	// parse document
	return parseDocument(r, docID)
}

func (storage *storageStupid) init() {
	storageInit()
}

func (storage *storageStupid) load(filename string) {
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

				pageIndex := parseFromFile(dir+"/"+file.Name(), docID)
				mergePageIndex(pageIndex, docID)
			} else {
				log.Println("Recursive is not supported")
			}
		}

		log.Println("Indexing directory, done")

		// debugPrintInvertedTable()
		storage.store(filename)
	}(storage.folderName)
}

func (storage *storageStupid) store(filename string) {
	serializing(filename, invertedIndex.data)
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
