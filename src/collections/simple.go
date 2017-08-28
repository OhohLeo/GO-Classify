package collections

import (
	"encoding/json"
)

func BuildSimple() Build {
	return func() Collection { return new(Simple) }
}

type Simple struct {
}

func (s *Simple) GetRef() Ref {
	return SIMPLE
}

func (s *Simple) Validate(id string, decoder *json.Decoder) error {
	return nil
}
