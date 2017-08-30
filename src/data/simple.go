package data

import (
	"time"
)

type Simple struct {
	Name string    `json:"name"`
	Time time.Time `json:"time"`
}

func (s *Simple) GetName() string {
	return s.Name
}

func (s *Simple) GetRef() Ref {
	return SIMPLE
}
