package main

import "github.com/go-ego/gse"

type segmentation interface {
	init()
	getSegmentedText(text []byte) []string
}

// gse

type segmentationGSE struct {
	segmenter gse.Segmenter
}

func (s *segmentationGSE) init() {
	s.segmenter.LoadDict()
}

func (s *segmentationGSE) getSegmentedText(text []byte) []string {
	segments := s.segmenter.Segment(text)

	return gse.ToSlice(segments, false)
}
