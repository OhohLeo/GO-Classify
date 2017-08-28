package collections

import (
	"encoding/json"
	//"fmt"
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

func NewMovieFromData(decoder *json.Decoder) (movie *Movie, err error) {
	err = decoder.Decode(&movie)
	return
}

func (m *Movie) Init() {
	m.Type = m.GetType()
}

func (m Movie) GetType() string {
	return "movie"
}

func (m *Movie) GetId() string {
	return m.Id
}

func (m *Movie) Update(decoder *json.Decoder) (err error) {
	err = decoder.Decode(&m)
	return
}

type Movies struct {
	Collection
}

// CreateItem create new movie item
func (m *Movies) CreateItem() *Movie {
	return new(Movie)
}

func (m *Movies) Validate(id string, decoder *json.Decoder) error {

	// Try to parse data received
	movie, err := NewMovieFromData(decoder)
	if err != nil {
		return err
	}

	// Force id
	movie.ItemGeneric.Id = id

	return nil
}
