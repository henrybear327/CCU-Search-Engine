package main

// Segmentation is the interface that difines Chinese segmentation
type Segmentation interface {
	init()
	getSegmentedText(text []byte) []string
}
