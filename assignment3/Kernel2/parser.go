package main

import (
	"log"
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

func parsePage(docID int, segmentedBody []string) map[string]*docTermData {
	result := make(map[string]*docTermData)

	pos := 0
	for _, token := range segmentedBody {
		if token, ok := tokenSanitizer(token); ok {
			log.Println(token)

			data, exist := result[token]
			if exist == false {
				data = &docTermData{docID: docID}
				result[token] = data
			}

			data.positions = append(data.positions, pos)
			pos++
		}
	}

	return result
}
