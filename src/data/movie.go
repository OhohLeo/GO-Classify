package data

import (
	"time"
)

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

func (m Movie) GetRef() Ref {
	return MOVIE
}

func (m *Movie) GetName() string {
	return m.Name
}
