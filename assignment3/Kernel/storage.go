package main

// Storage of the kernel
type Storage interface {
	/* Restores the inverted index hash map in the memory */
	init()
}
