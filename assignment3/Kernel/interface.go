package main

import "fmt"

func ui() {
	// fmt.Println(len(index), "words indexed in", len(indexedFiles), "files")
	for {
		fmt.Print("Search: ")
		var word string
		wc, _ := fmt.Scanln(&word)
		if wc == 0 {
			return
		}

		dl := index[word]
		switch len(dl) {
		case 0:
			fmt.Println("No match")
		case 1:
			fmt.Println("One match:")
		default:
			fmt.Println(len(dl), "matches:")
		}

		for key := range dl {
			fmt.Println("\t", indexedFiles[key].filename)
		}
	}
}
