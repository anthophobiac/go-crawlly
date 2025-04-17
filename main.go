package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly"
)

type HNItem struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

func main() {
	items, err := scrapeHackerNews()
	if err != nil {
		log.Fatal("Scraping failed:", err)
	}

	err = saveAsJSON("hackernews.json", items)
	if err != nil {
		log.Fatal("Failed to save JSON:", err)
	}

	fmt.Printf("✅  Scraped %d stories and saved to hackernews.json\n", len(items))
}

func scrapeHackerNews() ([]HNItem, error) {
	var items []HNItem

	c := colly.NewCollector(
		colly.AllowedDomains("news.ycombinator.com"),
		colly.UserAgent("Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"),
	)

	c.OnHTML("tr.athing", func(e *colly.HTMLElement) {
		title := e.ChildText(".titleline a")
		link := e.ChildAttr(".titleline a", "href")

		if title != "" && link != "" {
			items = append(items, HNItem{Title: title, Link: link})
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Request to %s failed: %v", r.Request.URL, err)
	})

	fmt.Println("Scraping Hacker News front page...")
	err := c.Visit("https://news.ycombinator.com/")
	if err != nil {
		return nil, err
	}

	return items, nil
}

func saveAsJSON(filename string, items []HNItem) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	return encoder.Encode(items)
}
