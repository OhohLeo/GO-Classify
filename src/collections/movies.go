package collections

import (
	"github.com/ohohleo/classify/websites"
)

type Movie struct {
	Status int
	Match  websites.Movie
	Founds []websites.Movie
}

type Movies struct {
	Collection
	movies            []Movie
	searchSubtitles   bool
	subtitleLanguages []string
	subtitles         map[string]websites.Website
}

// GetType returns the type of collection
func (m *Movies) GetType() string {
	return "movies"
}
