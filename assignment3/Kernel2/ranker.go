package main

type rankElement struct {
	termIdx   int
	positions []int
}

func rank(queries []string, matchedTermData []*termData) []int {
	// merge all docID results
	mapping := make(map[int][]*rankElement)
	for i, td := range matchedTermData {
		if td != nil {
			for _, rec := range td.documents {
				_, exist := mapping[rec.docID]
				if exist == false {
					mapping[rec.docID] = make([]*rankElement, 0)
				}
				mapping[rec.docID] = append(mapping[rec.docID], &rankElement{termIdx: i, positions: rec.positions})
			}
		}
	}

	// rank them
	var results []int
	for key, _ := range mapping {
		results = append(results, key)
	}

	return results // docID
}
