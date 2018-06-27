package main

import (
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

// perform segmentation
func (s *segmenter) getSegmentedText(text string) []string {
	byteString := []byte(text)
	segments := s.segmenter.Segment(byteString)

	return gse.ToSlice(segments, false)
}
