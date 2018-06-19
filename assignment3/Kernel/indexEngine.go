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
			termNode := &termNode{Total: 0, DocCount: 0, Data: make(map[int][]int)}
			invertedIndex.data[key] = termNode
		}
		invertedIndex.data[key].Total += len(value)
		invertedIndex.data[key].DocCount++
		invertedIndex.data[key].Data[docID] = value
	}
}
