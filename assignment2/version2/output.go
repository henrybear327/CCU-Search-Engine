package main

import (
	"bufio"
	"fmt"
	"os"
)

// OutputSeedingSites is a function that writes seeding sites to a file
func outputSeedingSites(seedingSites []string, conf *config) {
	f, err := os.Create(conf.Output.Seedfile)
	check(err)
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, rec := range seedingSites {
		check(err)
		fmt.Fprintln(w, rec)
	}
	w.Flush() // Don't forget to flush!
}
