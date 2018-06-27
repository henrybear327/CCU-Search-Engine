package main

import "sync"

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
	documents            []docTermData
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
func (i *indexer) merge(result map[string]*termData) {
	i.invertedTableLock.Lock()
	defer i.invertedTableLock.Unlock()
}

// insert
func (i *indexer) insert(title, body, url string) {
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
}

// query
func (i *indexer) query(query string) []*termData {
	seg.Lock()
	// segmentedQuery := seg.getSegmentedText(query)
	seg.Unlock()

	var results []*termData
	return results
}
