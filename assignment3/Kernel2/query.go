package main

import "log"

// returns total result count, query result as array of strings
func queryByString(query string, from int, to int) (int, []string) {
	totalRecords := 0
	results := make([]string, 0)

	matchedTermData := idxer.query(query)
	if len(matchedTermData) == 0 {
		return totalRecords, results
	}

	// from, to bound check
	return totalRecords, results
}

// returns docID
func insertDocument(title, body, url string) int {
	docID := idxer.insert(title, body, url)

	if debug {
		log.Println("Debug print inverted table")
		idxer.printInvertedTable()

		log.Println("Debug print database")
		idxer.printDatabase()
	}

	return docID
}
