package main

import (
	"fmt"

	"github.com/mmcdole/gofeed"
)

func rssTest() {
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL("https://hnrss.org/frontpage")
	fmt.Println(feed.Title, feed.UpdatedParsed)

	for _, item := range feed.Items {
		fmt.Println("=========")
		fmt.Println("Title", item.Title)             // feedly title
		fmt.Println("Description", item.Description) // feedly "excerpt"
		fmt.Println("Link", item.Link)               // feedly "visit website"
		fmt.Println("=========")
	}
}
