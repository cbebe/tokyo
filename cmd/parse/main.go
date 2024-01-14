package main

import (
	"log"
	"os"
	"path"

	"github.com/cbebe/tokyo"
)

func main() {
	p, err := tokyo.GetPageCSV("data.csv")
	if err != nil {
		log.Fatal(err)
	}

	places := tokyo.CreatePageMap(p)
	entries, err := os.ReadDir("html")
	if err != nil {
		log.Fatal(err)
	}

	entryJSON := make(map[string]tokyo.PageJSON, len(entries))
	for _, e := range entries {
		base := e.Name()
		f, err := os.Open(path.Join("html", base))
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		value, err := tokyo.GetPageJSON(f)
		if err != nil {
			log.Fatal(err)
		}
		entryJSON[places[base]] = *value
	}

	if err := tokyo.SaveAsJSON(entryJSON, "data.json"); err != nil {
		log.Fatal(err)
	}
}
