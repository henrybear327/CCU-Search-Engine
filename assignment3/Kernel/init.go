package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

func parse() (string, int) {
	source := flag.String("source", "docs", "the source folder to index")
	port := flag.Int("port", 8001, "port to listen for requests")
	flag.Parse()

	return *source, *port
}

func loadDataFromFile() {

}

type requestMessage struct {
	Operation int    `json:"operation"`
	Payload   string `json:"payload"`
}

func (r *requestMessage) String() string {
	return "Operation: " + strconv.Itoa(r.Operation) + "\n" + "Payload: " + r.Payload
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

func handleConnection(incomingRequest net.Conn) {
	defer func() {
		fmt.Fprintf(incomingRequest, "GET / HTTP/1.0\r\n\r\n")
		incomingRequest.Close()
	}()

	debugPrintRequest(incomingRequest)

	// var request requestMessage
	// err := json.NewDecoder(incomingRequest).Decode(&request)
	// check("json.NewDecoder(incomingRequest).Decode(&request)", err)

	// log.Println("Incoming request message:\n" + request.String())
}

func listenOnPort(port int) {
	/* Listen for incoming messages */
	ln, err := net.Listen("tcp", fmt.Sprint(":"+strconv.Itoa(port)))
	check("net.Listen", err)
	defer ln.Close()

	log.Println("Listening on port", port)

	/* accept connection on port */
	/* not sure if looping infinetely on ln.Accept() is good idea */
	for {
		incomingRequest, err := ln.Accept()
		if err != nil {
			log.Println("Error received while listening.", err)
			continue
		}
		go handleConnection(incomingRequest)
	}
}
