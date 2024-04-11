package main

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

// Movie holds the data about a movie
type Movie struct {
	Title string
}

func main() {
	// Instantiate default collector
	c := colly.NewCollector()

	// Counter variable to keep track of the number of movies scraped
	var movieCount int

	// OnHTML callback for the movie title based on the CSS selector
	c.OnHTML("h3.ipc-title__text", func(e *colly.HTMLElement) {
		// Stop scraping when we've collected the top 10 movies
		if movieCount >= 10 {
			return
		}

		// Extract the movie title text
		movieTitle := e.Text

		// This checks if the title includes a ranking number like "1. " and removes it
		parts := strings.SplitN(movieTitle, " ", 2)
		if len(parts) > 1 {
			movieTitle = parts[1]
		}

		movie := Movie{
			Title: movieTitle,
		}

		// Increment the movie count
		movieCount++

		// Print the movie title
		fmt.Printf("Top %d Movie: %s\n", movieCount, movie.Title)
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping on IMDb Top 250 Movies page
	c.Visit("https://www.imdb.com/chart/top/")
}
