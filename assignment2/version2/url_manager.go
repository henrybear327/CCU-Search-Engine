package main

import (
	"container/list"

	"github.com/temoto/robotstxt"
)

// Manager is the heart of every seed website
type Manager struct {
	link     string
	robot    *robotstxt.RobotsData
	urlQueue *list.List
	urlSeen  map[string]bool
}

func (manager *Manager) preprocess(done chan bool) {
	// parse robots.txt
	manager.parseRobotsTxt()

	// parse sitemap.xml (prepare queue)
	manager.parseSiteMap()

	done <- true
}
