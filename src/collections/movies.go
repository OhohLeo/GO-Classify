package collections

import (
	"encoding/json"
	//"fmt"
	"time"
)

func BuildMovies() Build {
	return func() Collection { return new(Movies) }
}

// Generic movie format
type Movie struct {
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

func NewMovieFromData(decoder *json.Decoder) (movie *Movie, err error) {
	err = decoder.Decode(&movie)
	return
}

func (m *Movie) Init() {
}

func (m Movie) GetType() string {
	return "movie"
}

func (m *Movie) GetId() string {
	return "123"
}

func (m *Movie) Update(decoder *json.Decoder) (err error) {
	err = decoder.Decode(&m)
	return
}

type Movies struct {
}

func (m *Movies) GetRef() Ref {
	return MOVIES
}

// CreateItem create new movie item
func (m *Movies) CreateItem() *Movie {
	return new(Movie)
}

func (m *Movies) Validate(id string, decoder *json.Decoder) error {

	// Try to parse data received
	// movie, err := NewMovieFromData(decoder)
	// if err != nil {
	// 	return err
	// }

	// Force id
	// movie.ItemGeneric.Id = id

	return nil
}
