package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	shindans, err := GetShindans("https://shindanmaker.com/c/list?mode=hot")
	if err != nil {
		log.Println(err)
	}
	shindanData, err := json.MarshalIndent(shindans, "", "  ")
	if err != nil {
		log.Println(err)
	}

	jsonShindan := string(shindanData)

	fmt.Println("🔥 Top 10 Hot Shindans 🔥")
	fmt.Printf("%v", jsonShindan)
}

// ShindanObj is a map of shindan title and urls
type ShindanObj struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

// GetShindans retrieves corona stats in json form
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
