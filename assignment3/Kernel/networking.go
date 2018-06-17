package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type insertionRequestMessage struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	URL   string `json:"url"`
}

func (r *insertionRequestMessage) String() string {
	return "Title: " + r.Title + "\nBody: " + r.Body + "\nURL: " + r.URL
}

type searchRequestMessage struct {
	Query string `json:"query"`
}

func (r *searchRequestMessage) String() string {
	return "Query: " + r.Query
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

	// Echo input JSON payload
	output, err := json.Marshal(msg)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(output)

	// print log
	log.Println("One insertion request is received", msg.String())
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
	log.Println("One search request is received", msg.String())

	// Perform searching
	results := textSearch(msg.Query)
	fmt.Println(results.String())

	// return result
	output, err := json.Marshal(results)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(output)
}

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
