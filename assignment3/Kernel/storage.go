package main

import (
	"encoding/gob"
	"log"
	"os"
	"sync"
)

/* DS definitions */
type termNode struct {
	Total    int
	DocCount int
	Data     map[int][]int
}

func (t *termNode) init() {
	t.Total = 0
	t.DocCount = 0
	t.Data = make(map[int][]int)
}

type invertedIndexData struct {
	sync.RWMutex
	data map[string]*termNode
}

/* serializing / deserializing */
func serializing(filePath string, object map[string]*termNode) {
	log.Println("Serializing started")

	file, err := os.Create(filePath)
	check("serializing", err)
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(object)
	check("encoder.Encode", err)

	log.Println("Serializing completed")
}

func deserializing(filePath string, object map[string]*termNode) {
	log.Println("Deserializing started")

	file, err := os.Open(filePath)
	check("deserializing", err)
	defer file.Close()

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&object)
	check("decoder.Decode", err)

	log.Println("Deserializing completed")
}

/* storage engine V1 */
type storageV1 struct {
	nextDocID     int
	invertedIndex invertedIndexData
}

func (storage *storageV1) init() {
	storage.nextDocID = 0
	storage.invertedIndex.data = make(map[string]*termNode) // term, (docID, [positions])
}

func (storage *storageV1) load(filename string) {
	// load key-value pairs from disk
	deserializing(filename, storage.invertedIndex.data)
}

func (storage *storageV1) store(filename string) {
	// store key-value pairs to disk
}

func (storage *storageV1) getAllTerms() []string {
	storage.invertedIndex.RLock()
	defer storage.invertedIndex.RUnlock()

	var ret []string
	for key := range storage.invertedIndex.data {
		ret = append(ret, key)
	}
	return ret
}

func (storage *storageV1) getTermRecords(term string) *termNode {
	storage.invertedIndex.RLock()
	defer storage.invertedIndex.RUnlock()

	return storage.invertedIndex.data[term]
}

func (storage *storageV1) insertTermRecord(term string, docID int, positions []int) {
	storage.invertedIndex.Lock()
	defer storage.invertedIndex.Unlock()

	if storage.invertedIndex.data[term] == nil {
		termNode := &termNode{Total: 0, DocCount: 0, Data: make(map[int][]int)}
		storage.invertedIndex.data[term] = termNode
	}
	storage.invertedIndex.data[term].Total += len(positions)
	storage.invertedIndex.data[term].DocCount++
	storage.invertedIndex.data[term].Data[docID] = positions
}
