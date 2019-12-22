package data

import (
	"time"
)

type Generic struct {
	Name string    `json:"name"`
	Time time.Time `json:"time"`
}

func (s *Generic) GetName() string {
	return s.Name
}

func (s *Generic) GetRef() Ref {
	return GENERIC
}
