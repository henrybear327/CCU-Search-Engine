package main

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
	idxer.insert(title, body, url)

	return 0
}
