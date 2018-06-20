package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

func check(functionName string, err error) {
	if err != nil {
		log.Fatalln("error from", functionName, err)
	}
}

func parseCommandLine() (string, int, string, int) {
	source := flag.String("source", "docs", "the source folder to index")
	port := flag.Int("port", 8001, "port to listen for requests")
	gobFile := flag.String("gobFile", "index.dat", "The index file to load")
	debug := flag.Int("debug", 0, "debug mode")
	flag.Parse()

	return *source, *port, *gobFile, *debug
}

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
	allTerms := config.storage.getAllTerms()
	for _, term := range allTerms {
		records := config.storage.getTermRecords(term)
		fmt.Println("[term]", term, "count", records.Total, "doc count", records.DocCount)

		for docID, positions := range records.Data {
			fmt.Printf("docID %v = [", docID)
			for _, position := range positions {
				fmt.Printf("%v ", position)
			}
			fmt.Println("]")
		}
	}
}
