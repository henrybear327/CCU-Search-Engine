package main

import (
	"bufio"
	"io"
)

func parseDocument(r *bufio.Reader, docID int) {
	// read input
	var str []byte
	for {
		b, isPrefix, err := r.ReadLine()
		switch {
		case err == io.EOF: // end of file
			return
		case err != nil: // non-EOF error, GG
			check("r.ReadLine()", err)
		case isPrefix == true:
			str = append(str, b...)
		case isPrefix == false:
			if len(str) == 0 {
				str = b
			} else {
				str = append(str, b...)
			}

			// Split by space
			// for _, bword := range bytes.Fields(str) {
			// 	bword := bytes.Trim(bword, ".,-~?!\"'`;:()<>[]{}\\|/=_+*&^%$#@")
			// 	if len(bword) > 0 {
			// 		word := string(bword)
			// 		docs := index[word]

			// 		if len(docs) == 0 {
			// 			index[word] = make(map[int]bool)
			// 			docs = index[word]
			// 		}

			// 		_, ok := docs[docID]
			// 		if ok == false {
			// 			docs[docID] = true
			// 		}
			// 	}
			// }

			// Split by segmentation
			for _, token := range configuration.segmenter.getSegmentedText(str) {
				// log.Println(token)
				if len(token) > 0 {
					docs := invertedIndex[token]

					if len(docs) == 0 {
						invertedIndex[token] = make(map[int]bool)
						docs = invertedIndex[token]
					}

					_, ok := docs[docID]
					if ok == false {
						docs[docID] = true
					}
				}
			}

			str = make([]byte, 0)
		}
	}
}
