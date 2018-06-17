package main

import (
	"bufio"
	"io"
	"strings"
)

func parseDocument(r *bufio.Reader, docID int) {
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

			// Perform segmentation
			for _, token := range configuration.segmenter.getSegmentedText(str) {
				if strings.Contains(token, " ") == true {
					continue
				}

				if len(token) > 0 {
					// fmt.Println(token)
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
