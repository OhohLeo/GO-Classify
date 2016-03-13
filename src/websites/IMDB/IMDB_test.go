package IMDB

import (
	"fmt"
	"github.com/ohohleo/classify/requests"
	"testing"
)

func TestSearch(t *testing.T) {

	requests.New(2, false)
	imdb := new(IMDB)

	c := imdb.SearchMovies("Star+Wars")

	for {
		movie, ok := <-c
		if ok {
			fmt.Printf("movie: %+v\n", movie)
			continue
		}

		break
	}

	//imdb.getResource("tt0405393")
}
