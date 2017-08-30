package collections

import (
	"encoding/json"
)

func BuildMovies() Build {
	return func() Collection { return new(Movies) }
}

type Movies struct {
}

func (m *Movies) GetRef() Ref {
	return MOVIES
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
