package main

import (
	"fmt"
	"log"
)

// returns total result count, query result as array of strings
func queryByString(query string, from int, to int) (int, []searchRequestReturnMessage) {
	totalRecords := 0
	results := make([]searchRequestReturnMessage, 0)

	segmentedQuery, matchedTermData := idxer.query(query)
	if len(matchedTermData) == 0 {
		if debug {
			fmt.Println("No matched term!")
		}
		return totalRecords, results
	}

	// rank
	if debug {
		fmt.Println("================================")
		fmt.Println("Segmented query")
		for i, td := range matchedTermData {
			if td != nil {
				fmt.Println("Term [", segmentedQuery[i], "] matched")
				printTermData(td) // TODO: possible race issue
			} else {
				fmt.Println("Term [", segmentedQuery[i], "] not matched")
			}
		}
		fmt.Println("================================")
	}
	rankedDocIDs := rank(segmentedQuery, matchedTermData)

	// [from, to) bound check, 0-based
	if from >= len(rankedDocIDs) {
		return totalRecords, results
	} else if to >= len(rankedDocIDs) {
		to = len(rankedDocIDs)
	}
	if debug {
		fmt.Println("Ranked docs", len(rankedDocIDs), "from", from, "to", to)
	}

	// highlight
	results = highlight(rankedDocIDs[from:to], segmentedQuery)

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
