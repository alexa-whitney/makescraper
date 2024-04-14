package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

// Movie holds the data about a movie
type Movie struct {
	Title   string
	Year    string
	Runtime string
	Rating  string
}

func main() {
	// Instantiate default collector
	c := colly.NewCollector()

	// Regular expression to match the year and rating
	reYear := regexp.MustCompile(`\b(19\d{2}|20\d{2})\b`)
	reRating := regexp.MustCompile(`R|PG-13|PG|G|NC-17`)

	// Counter variable to keep track of the number of movies scraped
	var movieCount int

	// OnHTML callback for the movie information based on the CSS selector
	c.OnHTML("div.ipc-metadata-list-summary-item__tc", func(e *colly.HTMLElement) {
		if movieCount >= 10 {
			return // Only scrape top 10 movies
		}

		// Extract the movie title text and remove the leading rank number
		movieTitleWithRank := e.ChildText("h3.ipc-title__text")
		movieTitle := strings.TrimSpace(strings.TrimLeft(movieTitleWithRank, "1234567890. "))

		// Initialize movie with title
		movie := Movie{
			Title: movieTitle,
		}

		// Extract the year, runtime, and rating based on their order
		// These details are in separate <span> elements; each item can be selected with its own class
		// Here we are assuming that the first span contains the year, the second the runtime, and the third the rating
		// This part might need to be adjusted if the assumption is incorrect
		movieDetails := e.DOM.Find("div").First().Find("span")
		if movieDetails.Length() >= 3 {
			// Extract year from the text
			yearText := movieDetails.Eq(0).Text()
			if matches := reYear.FindStringSubmatch(yearText); len(matches) > 1 {
				movie.Year = matches[1]
			}

			// Extract runtime directly
			movie.Runtime = strings.TrimSpace(movieDetails.Eq(1).Text())

			// Extract rating from the text
			ratingText := movieDetails.Eq(2).Text()
			if matches := reRating.FindString(ratingText); matches != "" {
				movie.Rating = matches
			}
		}

		// Increment the movie count
		movieCount++

		// Print the movie information
		fmt.Printf("Top %d Movie: %s (%s), Runtime: %s, Rating: %s\n", movieCount, movie.Title, movie.Year, movie.Runtime, movie.Rating)
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping on IMDb Top 250 Movies page
	c.Visit("https://www.imdb.com/chart/top/")
}
