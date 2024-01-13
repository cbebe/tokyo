package tokyo

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/PuerkitoBio/goquery"
)

func DownloadToFile(href, dir string) error {
	file, err := os.Create(path.Join("html", path.Base(href)))
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
	// Request the HTML page.
	res, err := http.Get(href)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}
	return res.Body, nil
}

func GetDocument(href string) (*goquery.Document, error) {
	body, err := GetPageBody(href)
	defer body.Close()
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}
	return doc, nil
}
