package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

func debugPrintRequest(incomingRequest net.Conn) {
	timeoutDuration := 1 * time.Second
	incomingRequest.SetReadDeadline(time.Now().Add(timeoutDuration))

	bufReader := bufio.NewReader(incomingRequest)
	var str []byte
	for {
		b, isPrefix, err := bufReader.ReadLine()
		if err == io.EOF { // end of file
			return
		} else if e, ok := err.(net.Error); ok && e.Timeout() {
			return
		} else if err != nil { // non-EOF error, GG
			check("bufReader.ReadLine()", err)
		} else if isPrefix == true {
			str = append(str, b...)
		} else if isPrefix == false {
			if len(str) == 0 {
				str = b
			} else {
				str = append(str, b...)
			}

			fmt.Println(string(str))

			str = make([]byte, 0)
		} else {
			log.Fatalln("WTF?")
		}
	}
}

func debugPrintInvertedTable() {
	for term, records := range invertedIndex.data {
		fmt.Println("term", term)
		for docID, positions := range records {
			fmt.Printf("docID %v = [", docID)
			for _, position := range positions {
				fmt.Printf("%v ", position)
			}
			fmt.Println("]")
		}
	}
}
