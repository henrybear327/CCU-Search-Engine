package main

func indexFromJSON(payload string) {

}

func mergePageIndex(pageIndex map[string][]int, docID int) {
	// merge index
	for key, value := range pageIndex {
		// fmt.Println(key, value)
		config.storage.insertTermRecord(key, docID, value)
	}
}
