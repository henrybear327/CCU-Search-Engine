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

func createFolderIfNotExist(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0777)
	}
}

func saveHTMLFileFromString(foldername, filename, pageSource string) {
	createFolderIfNotExist(conf.Output.PageSourcePath)
	path := conf.Output.PageSourcePath + "/" + foldername
	createFolderIfNotExist(path)

	f, err := os.Create(path + "/" + filename)
	check(err)
	defer f.Close()

	w := bufio.NewWriter(f)
	fmt.Fprintln(w, pageSource)
	w.Flush() // Don't forget to flush!
}
