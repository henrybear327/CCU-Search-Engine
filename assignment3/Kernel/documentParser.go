package main

import (
	"bufio"
	"io"
	"strings"
)

func tokenSanitizer(token string) (string, bool) {
	if strings.Contains(token, " ") == true {
		return "", false
	}

	if len(token) == 0 {
		return "", false
	}

	// TODO: regex strip special characters

	return token, true
}

func parseDocument(r *bufio.Reader, docID int) map[string][]int {
	// term: [pos]
	pageIndex := make(map[string][]int)
	position := 0

	var str []byte
	for {
		b, isPrefix, err := r.ReadLine()
		switch {
		case err == io.EOF: // end of file
			// for key, value := range pageIndex {
			// 	fmt.Println("[parseDocument]", key, value)
			// }
			return pageIndex
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
				if token, ok := tokenSanitizer(token); ok {
					// fmt.Println(token)

					pageIndex[token] = append(pageIndex[token], position)
					position++
				}
			}

			str = make([]byte, 0)
		}
	}
}
