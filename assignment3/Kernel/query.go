package main

import "fmt"

// SearchResult is the query result
type SearchResult struct {
	Count   int      `json:"count"`
	Results []string `json:"results"`
}

func (results *SearchResult) String() string {
	str := ""

	switch results.Count {
	case 0:
		str = fmt.Sprintln("No match")
	case 1:
		str = fmt.Sprintln("One match:")
	default:
		str = fmt.Sprintln(results.Count, "matches:")
	}
	for _, res := range results.Results {
		str += fmt.Sprintln("\t", res)
	}

	return str
}

func textSearch(query string) *SearchResult {
	dl := index[query]

	var results SearchResult
	results.Count = len(dl)
	for key := range dl {
		// fmt.Println("\t", indexedFiles[key].filename)
		results.Results = append(results.Results, indexedFiles[key].filename)
	}

	return &results
}

func searchUI() {
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
		fmt.Println(results.String())
	}
}
