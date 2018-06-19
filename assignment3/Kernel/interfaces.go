package main

// Segmentation defines Chinese segmentation engine
type Segmentation interface {
	init()
	getSegmentedText(text []byte) []string
}

// Storage of the kernel
type Storage interface {
	init()
	load() /* Restores the inverted index hash map in the memory */
}

func storageInit() {
	invertedIndex.data = make(map[string]map[int][]int) // term, (docID, [positions])
	indexedFiles.data = make(map[int]document)
}
