package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/cbebe/tokyo"
)

const ROOT = "https://www.gpsmycity.com"
const TOKYO = ROOT + "/gps-tour-guides/tokyo-405.html"

func ScrapeTokyo() {
	if err := os.MkdirAll("html", fs.ModePerm); err != nil {
		log.Fatal(err)
	}
	doc, err := tokyo.GetDocument(TOKYO)
	if err != nil {
		log.Fatal(err)
	}

	sb := strings.Builder{}

	// Find the review items
	doc.Find(".sfg_list.tbl").First().Children().Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Text())
		if strings.HasPrefix(title, "Custom Walk") {
			return
		}
		href, ok := s.Find("a").Attr("href")
		if !ok || (!strings.HasPrefix(href, "/tours") && !strings.HasPrefix(href, "/discovery")) {
			return
		}

		sb.WriteString("\"")
		sb.WriteString(title)
		sb.WriteString("\"")
		sb.WriteString(",")
		sb.WriteString(ROOT + href)
		sb.WriteString("\n")
		err := tokyo.DownloadToFile(ROOT+href, "html")
		if err != nil {
			log.Fatal(err)
		}
	})

	if err := os.WriteFile("data.csv", []byte(sb.String()), 0o644); err != nil {
		fmt.Println(sb.String())
		log.Fatal(err)
	}
}

func main() {
	ScrapeTokyo()
}
