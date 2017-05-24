package TMDB

import (
	"fmt"
	"testing"
)

func TestSearch(t *testing.T) {

	tmdb := New()

	tmdb.SetConfig(map[string]string{
		"api_key":     "80e1027731e519ce1f4faacf21ecc13a",
		"poster_path": "http://image.tmdb.org/t/p/original",
	})

	c := tmdb.Search("Medecin+de+Campagne")

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
