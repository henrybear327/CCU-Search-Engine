package main

func parsePage(docID int, segmentedBody []string) (map[string]*docTermData, int) {
	result := make(map[string]*docTermData)

	pos := 0
	for _, token := range segmentedBody {
		data, exist := result[token]
		if exist == false {
			data = &docTermData{docID: docID}
			result[token] = data
		}

		data.positions = append(data.positions, pos)
		pos++
	}

	return result, pos
}
