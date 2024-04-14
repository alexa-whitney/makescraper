package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

// Movie holds the data about a movie
type Movie struct {
	Title   string `json:"title"`
	Year    string `json:"year"`
	Runtime string `json:"runtime"`
	Rating  string `json:"rating"`
}

func main() {
	// Instantiate default collector
	c := colly.NewCollector()

	// Regular expression to match the year
	reYear := regexp.MustCompile(`\b(19\d{2}|20\d{2})\b`)

	// Slice to hold all movies
	var movies []Movie

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

		// Extract the year, runtime, and rating based on their order and class name
		// These details are in separate <span> elements; each item can be selected with its own class
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
			movie.Rating = strings.TrimSpace(movieDetails.Eq(2).Text())
		}

		// Append the movie to the slice
		movies = append(movies, movie)

		// Increment the movie count
		movieCount++
	})

	// After all movies are scraped, serialize them to JSON
	c.OnScraped(func(_ *colly.Response) {
		// Marshal the slice of movies to JSON
		jsonData, err := json.MarshalIndent(movies, "", "  ")
		if err != nil {
			fmt.Println("Error serializing JSON:", err)
			return
		}

		// Print JSON to stdout to validate it
		fmt.Println(string(jsonData))

		// Write JSON to a file
		err = os.WriteFile("output.json", jsonData, 0644)
		if err != nil {
			fmt.Println("Error writing JSON to file:", err)
			return
		}
	})

	// Start scraping
	c.Visit("https://www.imdb.com/chart/top/")
}
