package main

import (
	"fmt"
	service "paws-n-planes/internal/seatGuruScraper/service/scraper"
)

func main() {
	airlines := service.Scrape()

	for _, airline := range airlines {
		fmt.Println(*&airline)
	}
}
