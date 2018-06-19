package main

type configuration struct {
	segmenter Segmentation
	storage   Storage
}

func (config *configuration) init() {
	config.segmenter.init()

	config.storage.init() // init map
	config.storage.load() // load data
}
