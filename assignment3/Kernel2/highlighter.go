package main

import (
	"strings"
)

func highlight(rankedDocIDs []int, segmentedQuery []string) []searchRequestReturnMessage {
	var results []searchRequestReturnMessage

	idxer.databaseLock.RLock()
	defer idxer.databaseLock.RUnlock()
	for _, docID := range rankedDocIDs {
		record := idxer.database[docID]

		highlightedText := record.body
		for _, token := range segmentedQuery {
			newToken := "<span style=\"color:red;font-weight:bold\">" + token + "</span>"
			highlightedText = strings.Replace(highlightedText, token, newToken, -1)
		}
		results = append(results, searchRequestReturnMessage{
			Title: record.title,
			Body:  highlightedText,
			URL:   record.url,
		})
	}

	return results
}
