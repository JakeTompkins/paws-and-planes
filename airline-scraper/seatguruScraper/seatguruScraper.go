package seatguruScraper

import (
	"fmt"

	colly "github.com/gocolly/colly/v2"
)

const AIRLINE_LIST_URL = "https://seatguru.com/browseairlines"

type Airline struct {
	Name string
	Code string
	Url  string
}

func findAirlineLinks(e *colly.HTMLElement) {

}

func Scrape() {
	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit(AIRLINE_LIST_URL)
}
