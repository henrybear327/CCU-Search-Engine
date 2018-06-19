package main

func indexFromJSON(payload string) {

}

func mergePageIndex(pageIndex map[string][]int, docID int) {
	// merge index
	invertedIndex.Lock()
	defer invertedIndex.Unlock()
	for key, value := range pageIndex {
		// fmt.Println(key, value)
		if invertedIndex.data[key] == nil {
			invertedIndex.data[key] = make(map[int][]int)
		}
		invertedIndex.data[key][docID] = value
	}
}
