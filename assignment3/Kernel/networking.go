package main

import (
	"encoding/json"
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
	return "\nTitle: " + r.Title + "\n" + "Body: " + r.Body + "\n" + "URL: " + r.URL
}

func handleInsertionRequestConnection(w http.ResponseWriter, r *http.Request) {
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

	log.Println("One insertion request is received:", msg.String())
}

func listen(port int) {
	// define route
	http.HandleFunc("/insert", handleInsertionRequestConnection)

	// define port
	address := ":" + strconv.Itoa(port)

	// start
	log.Println("Starting server on address", address)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		panic(err)
	}
}
