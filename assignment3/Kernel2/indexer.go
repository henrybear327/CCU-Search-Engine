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

	wordCount int
}

// inverted table
type docTermData struct {
	docID     int
	positions []int
}

type termData struct {
	totalOccurrenceCount int
	documents            []*docTermData
}

type indexer struct {
	docID     int
	totalDocs int
	docIDLock sync.RWMutex

	invertedTable     map[string]*termData
	invertedTableLock sync.RWMutex

	database     map[int]*document
	databaseLock sync.RWMutex
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

	i.totalDocs++
	return ret
}

func (i *indexer) getCurrentDocID() int {
	i.docIDLock.RLock()
	defer i.docIDLock.RUnlock()

	return i.docID
}

// insert
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

func (i *indexer) insert(title, body, url string) int {
	seg.Lock()
	segmentedBody := seg.getSegmentedText(body)
	seg.Unlock()

	docID := i.getNextDocID()

	// to inverted table
	parsed, wordCount := parsePage(docID, segmentedBody)
	i.merge(parsed)

	// to database
	newDocument := &document{title, body, url, wordCount}
	i.databaseLock.Lock()
	i.database[docID] = newDocument
	i.databaseLock.Unlock()

	return docID
}

// query
func (i *indexer) query(query string) ([]string, []*termData) {
	seg.Lock()
	segmentedQuery := seg.getSegmentedText(query)
	seg.Unlock()

	i.invertedTableLock.Lock()
	defer i.invertedTableLock.Unlock()

	results := make([]*termData, 0)
	for _, term := range segmentedQuery {
		if _, ok := i.invertedTable[term]; ok {
			results = append(results, i.invertedTable[term])
		} else {
			results = append(results, nil)
		}
	}

	return segmentedQuery, results
}

// debug
func printTermData(data *termData) {
	for _, doc := range data.documents {
		fmt.Printf("\tdocID = %v (", doc.docID)
		for i, pos := range doc.positions {
			if i == 0 {
				fmt.Printf("")
			} else {
				fmt.Printf(", ")
			}
			fmt.Printf("%v", pos)
		}
		fmt.Println(")")
	}
}

func (i *indexer) printInvertedTable() {
	i.invertedTableLock.RLock()
	defer i.invertedTableLock.RUnlock()

	fmt.Println("=========================================")
	for key, value := range i.invertedTable {
		fmt.Println("[Key]", key, "\n[Total occurrence count]", value.totalOccurrenceCount)
		printTermData(value)
	}
	fmt.Println("=========================================")
}

func (i *indexer) printDatabase() {
	i.databaseLock.RLock()
	defer i.databaseLock.RUnlock()

	fmt.Println("=========================================")
	for key, value := range i.database {
		fmt.Println("[DocID]", key)
		fmt.Println("\t[Word count]", value.wordCount)
		fmt.Println("\t[Title]", value.title)
		fmt.Println("\t[Body]", value.body)
		fmt.Println("\t[URL]", value.url)
	}
	fmt.Println("=========================================")
}
