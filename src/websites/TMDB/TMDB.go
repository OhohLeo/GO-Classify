package TMDB

import (
	"fmt"
	"github.com/ohohleo/classify/data"
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

func (t *TMDB) GetRef() websites.Ref {
	return websites.TMDB
}

// Launch a search request through the3 IMDB Api
func (t *TMDB) Search(input string) chan data.Data {

	c := make(chan data.Data)

	go func() {

		currentPage := -1
		totalPages := 0

		// Get all pages
		for currentPage <= totalPages {

			options := make(map[string]string)

			if currentPage > 0 {
				options["page"] = strconv.Itoa(currentPage)
			}

			results := &api.MovieSearchResults{
				Page: 0,
				Results: []api.MovieShort{
					api.MovieShort{
						ID:            1,
						OriginalTitle: "original title",
						Title:         "title",
						ReleaseDate:   "2010-01-02",
						PosterPath:    "path",
					},
					api.MovieShort{
						ID:            2,
						OriginalTitle: "original title 2",
						Title:         "title 2",
						ReleaseDate:   "2011-01-02",
						PosterPath:    "path2",
					},
				},
				TotalPages:   0,
				TotalResults: 2,
			}

			// , err := t.api.SearchMovie(input, options)
			// if err != nil {
			// 	close(c)
			// 	return
			// }

			currentPage = results.Page
			totalPages = results.TotalPages

			for _, d := range results.Results {

				fmt.Printf("movie: %+v\n", d)

				release, _ := time.Parse("2006-01-02", d.ReleaseDate)

				movie := &data.Movie{
					Name:     d.OriginalTitle,
					Image:    t.posterPath + d.PosterPath,
					Released: release,
				}

				// movie.ItemGeneric.Id = t.GetName() + "_" + strconv.Itoa(d.ID)

				c <- movie
			}

			currentPage++
		}

		close(c)
	}()

	return c
}
