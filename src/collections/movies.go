package collections

import (
	"time"
)

// Generic movie format
type Movie struct {
	ItemGeneric
	Name        string    `json:"name"`
	Url         string    `json:"url"`
	Released    time.Time `json:"released"`
	Duration    int       `json:"duration"`
	Image       string    `json:"image"`
	Description string    `json:"description"`
	Directors   []string  `json:"directors"`
	Cast        []string  `json:"cast"`
	Genres      []string  `json:"genres"`
}

func (m *Movie) Init() {
	m.Type = m.GetType()
}

// GetType returns the type of collection
func (m Movie) GetType() string {
	return "movie"
}

type MovieItem struct {
	Item
	Match Movie
}

type Movies struct {
	Collection
	movies map[string]*MovieItem
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
