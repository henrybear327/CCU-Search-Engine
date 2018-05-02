package main

import "github.com/temoto/robotstxt"
import "container/list"

// Manager is the heart of every seed website
type Manager struct {
	link     string
	robot    *robotstxt.RobotsData
	urlQueue *list.List
	urlSeen  map[string]bool
}
