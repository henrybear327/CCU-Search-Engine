package main

import (
	"bufio"
	"fmt"
	"os"
)

func outputSeedingSites(seedingSites []string, seedingSitesOption []string) {
	f, err := os.Create(conf.Output.Seedfile)
	check(err)
	defer f.Close()

	w := bufio.NewWriter(f)
	for i, rec := range seedingSites {
		check(err)
		fmt.Fprintln(w, rec, seedingSitesOption[i])
	}
	w.Flush() // Don't forget to flush!
}

func saveHTMLFileFromString(foldername, filename, pageSource string) {
	path := conf.Output.PageSourcePath + "/" + foldername
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0777)
	}

	f, err := os.Create(path + "/" + filename)
	check(err)
	defer f.Close()

	w := bufio.NewWriter(f)
	fmt.Fprintln(w, pageSource)
	w.Flush() // Don't forget to flush!
}
