package collections

import (
	"encoding/json"
)

func BuildSimple() Build {
	return Build{
		ForceCreate: func() Collection {
			return new(Simple)
		},
		Create: func(json.RawMessage, json.RawMessage) (Collection, error) {
			return new(Simple), nil
		},
	}
}

type Simple struct {
}

func (s *Simple) GetRef() Ref {
	return SIMPLE
}

func (r *Simple) Check(config json.RawMessage) error {
	return nil
}

func (s *Simple) Validate(id string, decoder *json.Decoder) error {
	return nil
}
