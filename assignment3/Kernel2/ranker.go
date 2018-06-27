package main

import (
	"sort"
)

type rankElement struct {
	termIdx   int
	positions []int
}

type rankScore struct {
	docID int
	score float64
}

type rankScoreArray []rankScore

func (r rankScoreArray) Len() int           { return len(r) }
func (r rankScoreArray) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r rankScoreArray) Less(i, j int) bool { return r[i].score > r[j].score }

func getTFIDF() float64 {
	return 0.0
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
	var intermediate []rankScore
	for id, _ := range mapping {
		intermediate = append(intermediate, rankScore{docID: id, score: float64(id)})
	}
	sort.Sort(rankScoreArray(intermediate))
	// fmt.Println(intermediate)

	var results []int
	for _, val := range intermediate {
		results = append(results, val.docID)
	}

	return results // docID
}
