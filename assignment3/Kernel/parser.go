package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
)

func indexFile(filename string, docID int) {
	// open file
	f, err := os.Open(filename)
	check("os.Open", err)
	defer f.Close()

	// register new file
	newDocument := document{filename}
	indexedFiles[docID] = newDocument

	// scan lines, one by one
	r := bufio.NewReader(f)
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

			for _, bword := range bytes.Fields(b) {
				bword := bytes.Trim(bword, ".,-~?!\"'`;:()<>[]{}\\|/=_+*&^%$#@")
				if len(bword) > 0 {
					word := string(bword)
					docs := index[word]

					if len(docs) == 0 {
						index[word] = make(map[int]bool)
						docs = index[word]
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
