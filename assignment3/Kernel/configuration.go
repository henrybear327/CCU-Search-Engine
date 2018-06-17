package main

import "github.com/go-ego/gse"

// Option is the global conf struct
type Option struct {
	segmenter Segmentation
	storage   Storage
}

func (option *Option) init() {
	option.segmenter.init()
	option.storage.init()
}

/* Start Storage */
// debug purpose, load data straight from files in the folder
type storageInitFromFolder struct {
	folderName string
}

func (storage *storageInitFromFolder) init() {
	go indexFromDirectory(storage.folderName)
}

/* End Storage */

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

/* End Segmentation */
