package TMDB

import (
	"fmt"
	"github.com/ohohleo/classify/collections"
	"github.com/ohohleo/classify/websites"
	api "github.com/ryanbradynd05/go-tmdb"
	"strconv"
	"time"
)

type TMDB struct {
	api        *api.TMDb
	apiKey     string
	posterPath string
}

func New() *TMDB {
	return &TMDB{}
}

func (t *TMDB) SetConfig(config map[string]string) bool {

	fmt.Printf("SetConfig %+v\n", config)

	// Get API key
	apiKey, ok := config["api_key"]
	if ok {
		t.apiKey = apiKey
	}

	// Get Poster Path
	posterPath, ok := config["poster_path"]
	if ok {
		t.posterPath = posterPath
	}

	t.api = api.Init(apiKey)

	return true
}

func (t *TMDB) GetName() string {
	return "TMDB"
}

// Launch a search request through the IMDB Api
func (t *TMDB) Search(input string) chan websites.Data {

	c := make(chan websites.Data)

	go func() {

		currentPage := -1
		totalPages := 0

		// Get all pages
		for currentPage <= totalPages {

			options := make(map[string]string)

			if currentPage > 0 {
				options["page"] = strconv.Itoa(currentPage)
			}

			results, err := t.api.SearchMovie(input, options)
			if err != nil {
				close(c)
				return
			}

			currentPage = results.Page
			totalPages = results.TotalPages

			for _, data := range results.Results {
				//fmt.Printf("movie: %+v\n", data)

				release, _ := time.Parse("2006-01-02", data.ReleaseDate)

				movie := &collections.Movie{
					Name:     data.OriginalTitle,
					Image:    t.posterPath + data.PosterPath,
					Released: release,
				}

				movie.Init()

				c <- movie
			}

			currentPage++
		}

		close(c)
	}()

	return c
}
