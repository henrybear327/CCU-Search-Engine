package main

import (
	"encoding/gob"
	"log"
	"os"
	"sync"
)

type invertedIndexData struct {
	sync.RWMutex
	data map[string]map[int][]int
}

func storageInit() {
	invertedIndex.data = make(map[string]map[int][]int) // term, (docID, [positions])
}

type storageV1 struct {
}

func (storage *storageV1) init() {
	storageInit()
}

func (storage *storageV1) load(filename string) {
	// load key-value pairs from disk
	deserializing(filename, invertedIndex.data)
}

func (storage *storageV1) store(filename string) {
	// store key-value pairs to disk
}

func serializing(filePath string, object map[string]map[int][]int) {
	log.Println("Serializing started")

	file, err := os.Create(filePath)
	check("serializing", err)
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(object)
	check("encoder.Encode", err)

	log.Println("Serializing completed")
}

func deserializing(filePath string, object map[string]map[int][]int) {
	log.Println("Deserializing started")

	file, err := os.Open(filePath)
	check("deserializing", err)
	defer file.Close()

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&object)
	check("decoder.Decode", err)

	log.Println("Deserializing completed")
}
