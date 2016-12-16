package websites

import (
	"time"
)

type Website interface {
	GetName() string
	SetConfig(map[string]string) bool
	Search(string) chan Data
}

type Data interface {
	GetType() string
}

// Generic movie format
type Movie struct {
	Name        string
	Url         string
	Released    time.Time
	Duration    int
	Image       string
	Description string
	Directors   []string
	Cast        []string
	Genres      []string
}

func (m Movie) GetType() string {
	return "movie"
}
