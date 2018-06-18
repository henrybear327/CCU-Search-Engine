package main

// Storage of the kernel
type Storage interface {
	/* Restores the inverted index hash map in the memory */
	init()
	load()
}

func storageInit() {
	invertedIndex.data = make(map[string]map[int]bool)
	indexedFiles.data = make(map[int]document)
}
