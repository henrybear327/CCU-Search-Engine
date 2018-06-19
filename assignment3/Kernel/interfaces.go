package main

// Segmentation defines Chinese segmentation engine
type Segmentation interface {
	init()
	getSegmentedText(text []byte) []string
}

// Storage of the kernel
type Storage interface {
	init()                 /* Init data structures */
	load(filename string)  /* Load the inverted index from disk */
	store(filename string) /* Store the inverted index to disk */
}
