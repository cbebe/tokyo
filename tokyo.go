package tokyo

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const ROOT = "https://www.gpsmycity.com"
const TOKYO = ROOT + "/gps-tour-guides/tokyo-405.html"

type Page struct {
	Title string
	Href  string
}

type pageJSON struct {
	Distance float64    `json:"distance,omitempty"`
	Path     [][]any    `json:"path,string"`
	Pins     [][]string `json:"pins"`
}

type PageJSON struct {
	Distance float64     `json:"distance,omitempty"`
	Path     [][]float64 `json:"path"`
	Pins     [][]string  `json:"pins"`
}

func convertToFloats(pj *pageJSON) (*PageJSON, error) {
	value := new(PageJSON)
	value.Distance = pj.Distance
	value.Pins = pj.Pins
	value.Path = make([][]float64, len(pj.Path))
	for i, coord := range pj.Path {
		p := make([]float64, 2)
		for j, point := range coord {
			switch v := point.(type) {
			case float64:
				p[j] = v
			case string:
				var err error
				p[j], err = strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, err
				}
			default:
				return nil, fmt.Errorf("conversion from float to %T not supported", v)
			}
		}
		value.Path[i] = p
	}

	return value, nil
}

func SaveAsJSON(v any, file string) error {
	if b, err := json.Marshal(v); err != nil {
		return err
	} else if err := os.WriteFile(file, b, 0o644); err != nil {
		return err
	} else {
		return nil
	}
}

func GetPageJSON(r io.ReadCloser) (*PageJSON, error) {
	re, err := regexp.Compile("{.*}")
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}
	parsed := pageJSON{}
	done := false
	doc.Find("script:not([type]):not([src])").
		Each(func(i int, s *goquery.Selection) {
			if done {
				return
			}
			script := strings.TrimSpace(s.Text())
			if strings.HasPrefix(script, "jarr") {
				done = true
				m := re.Find([]byte(script))
				if m == nil {
					err = fmt.Errorf("pattern did not match!")
				} else {
					err = json.Unmarshal(m, &parsed)
				}
			}
		})

	return convertToFloats(&parsed)
}

func SaveToCSV(pages []Page) error {
	csvFile, err := os.Create("data.csv")
	if err != nil {
		return err
	}
	defer csvFile.Close()
	writer := csv.NewWriter(csvFile)
	for _, p := range pages {
		if err := writer.Write([]string{p.Title, p.Href}); err != nil {
			return err
		}
	}
	writer.Flush()
	return writer.Error()
}

func getDocument(href string) (*goquery.Document, error) {
	body, err := GetPageBody(href)
	defer body.Close()
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func ScrapeMainPageForLinks(root, href string) ([]Page, error) {
	doc, err := getDocument(href)
	if err != nil {
		return nil, err
	}
	pages := []Page{}
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
			pages = append(pages, Page{Title: title, Href: root + href})
		})

	return pages, nil
}

func GetPageCSV(file string) ([]Page, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	records, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}
	pages := make([]Page, 0, len(records))
	for _, r := range records {
		pages = append(pages, Page{r[0], r[1]})
	}
	return pages, nil
}

func CreatePageMap(pages []Page) map[string]string {
	places := make(map[string]string)
	for _, p := range pages {
		places[path.Base(p.Href)] = p.Title
	}
	return places
}

func DownloadPages(pages []Page, dir string) error {
	if err := os.MkdirAll(dir, fs.ModePerm); err != nil {
		return err
	}
	for _, p := range pages {
		if err := downloadToFile(p.Href, dir); err != nil {
			return err
		}
	}
	return nil
}

func downloadToFile(href, dir string) error {
	file, err := os.Create(path.Join(dir, path.Base(href)))
	if err != nil {
		return err
	}
	defer file.Close()
	tour, err := GetPageBody(href)
	if err != nil {
		return err
	}
	defer tour.Close()
	_, err = io.Copy(file, tour)
	return err
}

func GetPageBody(href string) (io.ReadCloser, error) {
	res, err := http.Get(href)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}
	return res.Body, nil
}
