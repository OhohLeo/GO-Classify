package collections

import (
	"github.com/ohohleo/classify/imports"
	"github.com/ohohleo/classify/requests"
	"github.com/ohohleo/classify/websites/IMDB"
	"log"
	"testing"
)

func TestMovies(t *testing.T) {

	requests.New(2, false)

	m := new(Movies)
	m.Register("imdb", new(IMDB.IMDB))
	c := m.OnInput(&imports.File{Path: "Looper"})

	for {
		if movie, ok := <-c; ok {
			log.Printf("%+v", movie.GetType())
			continue
		}
		break
	}
}
