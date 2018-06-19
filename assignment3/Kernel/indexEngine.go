package main

import "sync"

type invertedIndexData struct {
	sync.RWMutex
	data map[string]map[int][]int
}

type indexedFilesData struct {
	sync.RWMutex
	data map[int]document
}

type document struct {
	filename string
}

func indexFromJSON(payload string) {

}
