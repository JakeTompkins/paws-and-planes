package seatguruScraper

import (
	"fmt"
	"strings"

	goquery "github.com/PuerkitoBio/goquery"
	colly "github.com/gocolly/colly/v2"
)

const ROOT_URL = "https://seatguru.com"
const AIRLINE_LIST_URL = "https://seatguru.com/browseairlines"

type Airline struct {
	Name string
	Code string
	Url  string
}

func findCode(airline *Airline, dom *goquery.Selection) {
	code := dom.Find("span.ai-label").FilterFunction(func(i int, sel *goquery.Selection) bool {
		return strings.HasPrefix(sel.Text(), "Airline Code")
	}).First().Next().Text()

	airline.Code = code
}

func findUrl(airline *Airline, dom *goquery.Selection) {
	officialUrl := dom.Find("span.ai-label").FilterFunction(func(i int, sel *goquery.Selection) bool {
		return strings.HasPrefix(sel.Text(), "Website")
	}).First().Next().Children().First().Text()

	airline.Url = officialUrl
}

func getAirlineDetials(url string, airline *Airline, rootColly *colly.Collector) {
	c := rootColly.Clone()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnHTML("body", func(e *colly.HTMLElement) {
		dom := e.DOM

		findCode(airline, dom)
		findUrl(airline, dom)
	})

	c.Visit(url)
}

func Scrape() []*Airline {
	airlines := []*Airline{}

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnHTML("div.browseAirlines > ul > li > a[href]", func(e *colly.HTMLElement) {
		airline := &Airline{
			Name: e.Text,
		}

		url := fmt.Sprintf("%s%s", ROOT_URL, e.Attr("href"))

		getAirlineDetials(url, airline, c)

		airlines = append(airlines, airline)
	})

	c.Visit(AIRLINE_LIST_URL)

	return airlines
}
