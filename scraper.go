package main

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// ShindanObj is a map of shindan title and urls
type ShindanObj struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

// GetShindans retrieves shindan list in json form
func GetShindans(url string) (map[int]ShindanObj, error) {

	resp, err := http.Get(url)
	if err != nil {
		res := make(map[int]ShindanObj)
		return res, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return make(map[int]ShindanObj), err
	}

	shindanMap := make(map[int]ShindanObj)

	doc.Find(".list_title").Each(func(i int, s *goquery.Selection) {
		if i < 10 {
			link, _ := s.Attr("href")
			url := "https://shindanmaker.com" + link
			shindanMap[i+1] = ShindanObj{s.Text(), url}
		}
	})

	return shindanMap, nil
}
