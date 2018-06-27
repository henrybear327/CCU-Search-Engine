package main

import (
	"strings"
	"sync"

	"github.com/go-ego/gse"
)

// setup
type segmenter struct {
	sync.RWMutex
	segmenter gse.Segmenter
}

func (s *segmenter) init() {
	s.segmenter.LoadDict()
}

func (s *segmenter) tokenSanitizer(token string) (string, bool) {
	if strings.Contains(token, " ") == true {
		return "", false
	}

	if len(token) == 0 {
		return "", false
	}

	// TODO: regex strip special characters

	return token, true
}

// perform segmentation
func (s *segmenter) getSegmentedText(text string) []string {
	byteString := []byte(text)
	segments := s.segmenter.Segment(byteString)

	results := make([]string, 0)

	for _, token := range gse.ToSlice(segments, false) {
		if token, ok := s.tokenSanitizer(token); ok {
			// log.Println(token)
			results = append(results, token)
		}
	}

	return results
}
