package main

import (
	"bytes"
	"fmt"
	"io"
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

	sb := new(bytes.Buffer)

	// Find the review items
	doc.Find(".sfg_list.tbl").
		First().
		Children().
		Each(func(i int, s *goquery.Selection) {
			title := strings.TrimSpace(s.Text())
			if strings.HasPrefix(title, "Custom Walk") {
				return
			}
			href, ok := s.Find("a").Attr("href")
			if !ok || (!strings.HasPrefix(href, "/tours") &&
				!strings.HasPrefix(href, "/discovery")) {
				return
			}

			fmt.Fprintf(sb, "\"%s\",%s\n", title, ROOT+href)
			if err := tokyo.DownloadToFile(ROOT+href, "html"); err != nil {
				log.Fatal(err)
			}
		})

	if csv, err := os.Create("data.csv"); err != nil {
		fmt.Println(sb.String())
		log.Fatal(err)
	} else if _, err = io.Copy(csv, sb); err != nil {
		fmt.Println(sb.String())
		log.Fatal(err)
	}
}

func main() {
	ScrapeTokyo()
}
