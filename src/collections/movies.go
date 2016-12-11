package collections

import (
	"github.com/ohohleo/classify/websites"
)

type Movie struct {
	Item
	Match websites.Movie
}

type Movies struct {
	Collection
	movies map[string]*Movie
}

// GetType returns the type of collection
func (m *Movies) GetType() string {
	return "movies"
}

// CreateItem create new movie item
func (m *Movies) CreateItem() *Movie {
	return new(Movie)
}

func (m *Movies) Validate(movie *Movie) {

}
