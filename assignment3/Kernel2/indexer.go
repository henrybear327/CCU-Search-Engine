package main

import (
	"fmt"
	"sync"
)

// collection database
type document struct {
	title string
	body  string
	url   string
}

// inverted table
type docTermData struct {
	docID     int
	positions []int
}

type termData struct {
	documents            []*docTermData
	totalOccurrenceCount int
}

type indexer struct {
	docID     int
	docIDLock sync.RWMutex

	invertedTable     map[string]*termData
	invertedTableLock sync.RWMutex
	database          map[int]*document
	databaseLock      sync.RWMutex
}

// init
func (i *indexer) init() {
	i.docID = 0
	i.invertedTable = make(map[string]*termData)
	i.database = make(map[int]*document)
}

// helper
func (i *indexer) getNextDocID() int {
	i.docIDLock.Lock()
	defer i.docIDLock.Unlock()

	ret := i.docID
	i.docID++
	return ret
}

// merge
func (i *indexer) merge(result map[string]*docTermData) {
	i.invertedTableLock.Lock()
	defer i.invertedTableLock.Unlock()

	for key, value := range result {
		data, exist := i.invertedTable[key]
		if exist == false {
			data = &termData{}
			i.invertedTable[key] = data
		}
		data.documents = append(data.documents, value)
		data.totalOccurrenceCount += len(value.positions)
	}
}

// insert
func (i *indexer) insert(title, body, url string) int {
	seg.Lock()
	segmentedBody := seg.getSegmentedText(body)
	seg.Unlock()

	// to database
	docID := i.getNextDocID()
	newDocument := &document{title, body, url}
	i.databaseLock.Lock()
	i.database[docID] = newDocument
	i.databaseLock.Unlock()

	// to inverted table
	parsed := parsePage(docID, segmentedBody)
	i.merge(parsed)

	return docID
}

// query
func (i *indexer) query(query string) []*termData {
	seg.Lock()
	// segmentedQuery := seg.getSegmentedText(query)
	seg.Unlock()

	var results []*termData
	return results
}

// debug
func (i *indexer) printInvertedTable() {
	fmt.Println("=========================================")
	for key, value := range i.invertedTable {
		fmt.Println("[Key]", key, "\n[Total occurrence count]", value.totalOccurrenceCount)
		for _, doc := range value.documents {
			fmt.Println("docID = ", doc.docID)
			for _, pos := range doc.positions {
				fmt.Printf("%v ", pos)
			}
			fmt.Println()
		}
	}
	fmt.Println("=========================================")
}
