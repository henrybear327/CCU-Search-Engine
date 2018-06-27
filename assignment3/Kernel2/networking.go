package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// networking
func listen(port int) {
	// define route
	http.HandleFunc("/insert", handleInsertionRequest)
	http.HandleFunc("/search", handleSearchRequest)

	// define port
	address := ":" + strconv.Itoa(port)

	// start
	log.Println("Starting server on address", address)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		panic(err)
	}
}

// http query processing
type insertionRequestMessage struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	URL   string `json:"url"`
	DocID int    `json:"docID"`
}

func (r *insertionRequestMessage) String() string {
	return "Title: " + r.Title + "\nBody: " + r.Body + "\nURL: " + r.URL
}

type searchRequestMessage struct {
	Query string `json:"query"`
	From  int    `json:"from"`
	To    int    `json:"to"`
}

func (r *searchRequestMessage) String() string {
	return "Query: " + r.Query + "\nfrom: " + strconv.Itoa(r.From) + "\nTo:" + strconv.Itoa(r.To)
}

func handleInsertionRequest(w http.ResponseWriter, r *http.Request) {
	// Read body
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Unmarshal JSON
	var msg insertionRequestMessage
	err = json.Unmarshal(b, &msg)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// insert
	msg.DocID = insertDocument(msg.Title, msg.Body, msg.URL)

	// Echo input JSON payload
	output, err := json.Marshal(msg)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(output)

	// print log
	log.Println("One insertion request is received\n", msg.String())
}

func handleSearchRequest(w http.ResponseWriter, r *http.Request) {
	// Read body
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Unmarshal JSON
	var msg searchRequestMessage
	err = json.Unmarshal(b, &msg)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	log.Println("One search request is received\n", msg.String())

	// Perform searching
	totalRecords, results := queryByString(msg.Query, msg.From, msg.To)
	fmt.Println(totalRecords, results)

	// return result
	output, err := json.Marshal(results)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(output)
}
