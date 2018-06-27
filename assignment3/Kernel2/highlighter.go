package main

func highlight(rankedDocIDs []int, segmentedQuery []string) []searchRequestReturnMessage {
	var results []searchRequestReturnMessage

	idxer.databaseLock.RLock()
	defer idxer.databaseLock.RUnlock()
	for _, docID := range rankedDocIDs {
		record := idxer.database[docID]
		results = append(results, searchRequestReturnMessage{
			Title: record.title,
			Body:  record.body,
			URL:   record.url,
		})
	}

	return results
}
