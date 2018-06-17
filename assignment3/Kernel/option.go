package main

import "github.com/go-ego/gse"

// Option is the global conf struct
type Option struct {
	segmenter Segmentation
}

/* Start Segmentation */
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

// cppJieba
type segmentationJieba struct {
	segmenter gse.Segmenter
}

func (s *segmentationJieba) init() {
	// load dictionary

}

func (s *segmentationJieba) getSegmentedText(text []byte) []string {
	// call segmentator

	// return []string
	return make([]string, 0)
}

/* End Segmentation */
