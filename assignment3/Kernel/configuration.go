package main

type configuration struct {
	segmenter Segmentation
	storage   Storage
}

/* CRUCIAL: register all init functions here! */
func (config *configuration) init(gobFilename string) {
	config.segmenter.init() // load dict

	config.storage.init()            // make map
	config.storage.load(gobFilename) // load data
}
