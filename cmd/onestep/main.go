package main

import (
	"log"

	"github.com/cbebe/tokyo"
)

func main() {
	pages, err := tokyo.ScrapeMainPageForLinks(tokyo.ROOT, tokyo.TOKYO)
	if err != nil {
		log.Fatal(err)
	}

	entryJSON := make(map[string]tokyo.PageJSON, len(pages))
	for _, p := range pages {
		body, err := tokyo.GetPageBody(p.Href)
		if err != nil {
			log.Fatal(err)
		}
		defer body.Close()

		value, err := tokyo.GetPageJSON(body)
		entryJSON[p.Title] = *value
	}

	if err := tokyo.SaveAsJSON(entryJSON, "data.json"); err != nil {
		log.Fatal(err)
	}
}
