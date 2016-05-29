package collections

import (
	"github.com/ohohleo/classify/imports"
	"github.com/ohohleo/classify/websites"
	"log"
)

type Movie struct {
	CollectionData
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
