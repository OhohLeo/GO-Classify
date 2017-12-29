package collections

import (
	"encoding/json"
)

func BuildMovies() Build {
	return Build{
		Create: func(json.RawMessage, json.RawMessage) (Collection, error) {
			return new(Movies), nil
		},
	}
}

type Movies struct {
}

func (m *Movies) GetRef() Ref {
	return MOVIES
}

func (r *Movies) Check(config json.RawMessage) error {
	return nil
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
