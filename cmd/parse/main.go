package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func readCSV() (map[string]string, error) {
	csv, err := os.Open("data.csv")
	if err != nil {
		return nil, err
	}
	places := make(map[string]string)
	s := bufio.NewScanner(csv)
	for s.Scan() {
		s := strings.Split(s.Text(), ",")
		var name string
		json.Unmarshal([]byte(s[0]), &name)
		places[path.Base(s[1])] = name
	}
	return places, err
}

func main() {
	places, err := readCSV()
	if err != nil {
		log.Fatal(err)
	}
	entries, err := os.ReadDir("html")
	if err != nil {
		log.Fatal(err)
	}
	re, err := regexp.Compile("{.*}")
	if err != nil {
		log.Fatal(err)
	}

	entryJson := make(map[string]any, len(entries))
	for _, e := range entries {
		base := e.Name()
		f, err := os.Open(path.Join("html", base))
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		doc, err := goquery.NewDocumentFromReader(f)
		if err != nil {
			log.Fatal(err)
		}
		found := false
		doc.Find("script:not([type]):not([src])").
			Each(func(i int, s *goquery.Selection) {
				if found {
					return
				}
				script := strings.TrimSpace(s.Text())
				if strings.HasPrefix(script, "jarr") {
					var v map[string]any
					found = true
					m := re.Find([]byte(script))
					if m == nil {
						log.Fatal("pattern did not match!")
					}
					err := json.Unmarshal(m, &v)
					if err != nil {
						log.Fatal(err)
					}

					entryJson[places[base]] = v
				}
			})
	}
	b, err := json.Marshal(entryJson)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile("data.json", b, 0o644)
	if err != nil {
		log.Fatal(err)
	}
}
