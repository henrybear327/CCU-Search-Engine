package main

import (
	"fmt"
	"math"
	"sort"
)

type rankElement struct {
	termIdx   int
	positions []int
}

type rankScore struct {
	docID          int
	proximityScore int
	tfidf          float64
}

type rankScoreArray []rankScore

func (r rankScoreArray) Len() int      { return len(r) }
func (r rankScoreArray) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r rankScoreArray) Less(i, j int) bool {
	if r[i].proximityScore == r[j].proximityScore { // proximity score
		return r[i].tfidf > r[j].tfidf // tf-idf
	}
	return r[i].proximityScore < r[j].proximityScore

	/*
		// bug on >
		if r[i].proximityScore < r[j].proximityScore { // proximity score
			return true
		}
		return r[i].tfidf > r[j].tfidf // tf-idf
	*/
}

// Implement TF-IDF
func getTFIDF(docID int, queries []string, elements []*rankElement) float64 {
	idxer.databaseLock.RLock()
	defer idxer.databaseLock.RUnlock()
	idxer.invertedTableLock.RLock()
	defer idxer.invertedTableLock.RUnlock()

	result := 0.0
	for _, element := range elements {
		// tf
		// ln(1 + number of times term t appears in the document / total terms in the document)
		// tf := math.Log(1.0 + float64(len(element.positions))/float64(idxer.database[docID].wordCount))
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

func dfs(idx int, elements []*rankElement, acc int, bestSoFar *int) {
	if idx >= len(elements)-1 {
		if acc < *bestSoFar {
			*bestSoFar = acc
		}
		return
	}

	if acc > *bestSoFar {
		return
	}

	for _, target := range elements[idx].positions {
		// fmt.Println("idx", idx, "target", target)
		// search at elements[idx + 1]
		lb := -1
		mx := len(elements[idx+1].positions) //[lb, ub)
		ub := mx
		for ub-lb > 1 {
			mid := lb + (ub-lb)/2
			if elements[idx+1].positions[mid] <= target {
				lb = mid
			} else {
				ub = mid
			}
		}

		if ub == mx {
			break
		}
		// fmt.Println("ub", ub, "nxt", elements[idx+1].positions[ub])
		dfs(idx+1, elements, acc+(elements[idx+1].positions[ub]-target-1), bestSoFar)
	}
}

func getProximity(docID int, queries []string, elements []*rankElement, queryLen int) int {
	result := 100 // threshold as init value
	dfs(0, elements, 0, &result)

	// fmt.Println(queries)
	// for _, val := range elements {
	// 	fmt.Printf("%v ", queries[val.termIdx])
	// }
	// fmt.Println()

	return result + (queryLen - len(elements)) //  penalty = terms missed
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
		intermediate = append(intermediate, rankScore{
			docID:          id,
			proximityScore: getProximity(id, queries, mapping[id], len(queries)),
			tfidf:          getTFIDF(id, queries, mapping[id]),
		})
	}
	sort.Sort(rankScoreArray(intermediate))
	fmt.Println(intermediate)

	var results []int
	for _, val := range intermediate {
		results = append(results, val.docID)
	}

	return results // docID
}
