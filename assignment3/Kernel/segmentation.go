package main

// Segmentation is the interface that handles 斷詞
type Segmentation interface {
	init()
	getSegmentedText(text []byte) []string
}
