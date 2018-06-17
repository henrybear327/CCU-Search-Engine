package main

import "fmt"

// SearchResult is the query result
type SearchResult struct {
	Results []string `json:"results"`
}

func textSearch(query string) *SearchResult {
	dl := index[query]
	switch len(dl) {
	case 0:
		fmt.Println("No match")
	case 1:
		fmt.Println("One match:")
	default:
		fmt.Println(len(dl), "matches:")
	}

	var results SearchResult
	for key := range dl {
		// fmt.Println("\t", indexedFiles[key].filename)
		results.Results = append(results.Results, indexedFiles[key].filename)
	}

	return &results
}

func ui() {
	fmt.Println("query string length = 0 -> quit")
	for {
		// get query
		fmt.Print("Search: ")
		var word string
		wc, _ := fmt.Scanln(&word)
		if wc == 0 {
			return
		}

		results := textSearch(word)
		for _, res := range results.Results {
			fmt.Println("\t", res)
		}
	}
}
