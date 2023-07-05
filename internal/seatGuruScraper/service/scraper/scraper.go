package scraper

import (
	"fmt"
	"strings"
	"time"

	a "paws-n-planes/pkg/models/airline"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

const ROOT_URL = "https://seatguru.com"
const AIRLINE_LIST_URL = "https://seatguru.com/browseairlines"

func getAirlineDetials(url string, airline *a.Airline, rootColly *colly.Collector) {
	c := rootColly.Clone()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		if r.StatusCode != 200 {
			fmt.Println(r.StatusCode)
		}
	})

	c.OnHTML("body", func(e *colly.HTMLElement) {
		detailsContainer := e.DOM.Find("div.airlineBannerLargeRight").First()

		codeTitle := detailsContainer.Find("span").FilterFunction(func(i int, s *goquery.Selection) bool {
			return strings.Contains(s.Text(), "Code")
		})
		code := codeTitle.Next().Text()

		airline.Code = code

		website, exists := detailsContainer.Find("a[href]").First().Attr("href")

		if exists == true {
			airline.Url = website
		}
	})

	c.Visit(url)
}

func Scrape() []*a.Airline {
	airlines := []*a.Airline{}

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnHTML("div.browseAirlines > ul > li > a[href]", func(e *colly.HTMLElement) {
		airline := &a.Airline{
			Name: e.Text,
		}

		url := fmt.Sprintf("%s%s", ROOT_URL, e.Attr("href"))

		time.Sleep(1 * time.Second) // NOTE: Sometimes the site returns blank, this should fix that

		getAirlineDetials(url, airline, c)

		airlines = append(airlines, airline)
	})

	c.Visit(AIRLINE_LIST_URL)

	return airlines
}
