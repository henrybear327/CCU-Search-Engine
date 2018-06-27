package main

import (
	"math"
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

// TODO
// Implement TF-IDF
func getTFIDF(docID int, queries []string, elements []*rankElement) float64 {
	idxer.databaseLock.RLock()
	defer idxer.databaseLock.RUnlock()
	idxer.invertedTableLock.RLock()
	defer idxer.invertedTableLock.RUnlock()

	result := 0.0
	for _, element := range elements {
		// tf
		// number of times term t appears in the document / total terms in the document
		tf := float64(len(element.positions)) / float64(idxer.database[docID].wordCount)

		// idf
		// ln(total number of documents / number of documents with term t in it)
		totalNumberOfDocuments := float64(idxer.totalDocs)
		numberOfDocsWithTerm := float64(len(idxer.invertedTable[queries[element.termIdx]].documents))
		idf := math.Log(totalNumberOfDocuments / numberOfDocsWithTerm)

		tfidf := tf * idf

		result += tfidf
	}

	return result
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
		intermediate = append(intermediate, rankScore{docID: id, score: getTFIDF(id, queries, mapping[id])})
	}
	sort.Sort(rankScoreArray(intermediate))
	// fmt.Println(intermediate)

	var results []int
	for _, val := range intermediate {
		results = append(results, val.docID)
	}

	return results // docID
}
