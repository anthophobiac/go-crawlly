# Hacker News Scraper (Go + Colly)

A simple web scraper that collects story titles and URLs from the [Hacker News](https://news.ycombinator.com/) front page. 
Built using [Colly](https://github.com/gocolly/colly).

## Features

- Scrapes titles and links from the Hacker News front page
- Supports concurrent scraping for multiple pages.
- Saves data to a formatted `hackernews.json` file

## Getting Started

* Clone the repo
```bash
git clone https://github.com/anthophobiac/go-crawlly.git
cd go-crawlly
```
* Make sure you have Go installed and then run the following to install dependencies:
```bash
go mod tidy
```
* Run the scraper
```bash
go run main.go
```

You can also scrape multiple pages concurrently:
```bash
go run main.go --pages 3
```

After completion, you’ll see a `hackernews.json` file with the scraped stories.

### Example Output
```json
[
  {
    "title": "Show HN: I built a search engine for open source",
    "link": "https://example.com"
  },
  ...
]
```