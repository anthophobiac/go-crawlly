package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gocolly/colly"
)

type HNItem struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

func main() {
	pages := flag.Int("pages", 1, "Number of Hacker News pages to scrape concurrently")
	flag.Parse()

	if *pages < 1 {
		log.Fatal("Pages must be at least 1.")
	}

	items, err := scrapeHackerNews(*pages)
	if err != nil {
		log.Fatal("Scraping failed:", err)
	}

	if err := saveAsJSON("hackernews.json", items); err != nil {
		log.Fatal("Saving JSON failed:", err)
	}

	fmt.Printf("✅ Scraped %d items from %d page(s), saved to hackernews.json\n", len(items), *pages)
}

func scrapeHackerNews(pages int) ([]HNItem, error) {
	var (
		items []HNItem
		mutex sync.Mutex
	)

	c := colly.NewCollector(
		colly.AllowedDomains("news.ycombinator.com"),
		colly.Async(true),
		colly.UserAgent("Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"),
	)

	err := c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 4,
		RandomDelay: 1 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	c.OnHTML("tr.athing", func(e *colly.HTMLElement) {
		title := e.ChildText(".titleline a")
		link := e.ChildAttr(".titleline a", "href")

		if title != "" && link != "" {
			mutex.Lock()
			items = append(items, HNItem{Title: title, Link: link})
			mutex.Unlock()
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		pageURL := r.Request.URL.String()

		if r.StatusCode == 503 {
			fmt.Printf("⚠️  Service Unavailable for %s — retrying in 3s...\n", pageURL)
			time.Sleep(3 * time.Second)

			retryErr := r.Request.Retry()
			if retryErr != nil {
				fmt.Printf("Retry failed for %s: %v\n", pageURL, retryErr)
			} else {
				fmt.Printf("✅ Successfully retried %s\n", pageURL)
			}
			return
		}

		fmt.Printf("Failed to load %s: %v\n", pageURL, err)
	})

	for i := 1; i <= pages; i++ {
		url := "https://news.ycombinator.com/news?p=" + strconv.Itoa(i)
		err := c.Visit(url)
		if err != nil {
			return nil, err
		}
	}

	c.Wait()
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
