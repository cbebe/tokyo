package main

import (
	"fmt"
	"log"

	"github.com/cbebe/tokyo"
)

func main() {
	pages, err := tokyo.ScrapeMainPageForLinks(tokyo.ROOT, tokyo.TOKYO)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pages, len(pages))
	err = tokyo.SaveToCSV(pages)
	if err != nil {
		log.Fatal(err)
	}
	err = tokyo.DownloadPages(pages, "html")
	if err != nil {
		log.Fatal(err)
	}
}
