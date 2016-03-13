package collections

import (
	"github.com/ohohleo/classify/imports"
	"github.com/ohohleo/classify/websites"
	"log"
)

type Movie struct {
	Status int
	Import []imports.Data
	Match  websites.Movie
	Founds []websites.Movie
	// Exports  []Export
}

type Movies struct {
	movies   []Movie
	websites map[string]websites.Website

	searchSubtitles   bool
	subtitleLanguages []string
	subtitles         map[string]websites.Website
}

// GetType returns the type of collection
func (m *Movies) GetType() string {
	return "movies"
}

// Add new movie website
func (m *Movies) Register(name string, website websites.Website) {

	if m.websites == nil {
		m.websites = make(map[string]websites.Website)
	}

	m.websites[name] = website
}

// OnInput handle new data to classify
func (m *Movies) OnInput(input imports.Data) chan websites.Data {

	c := make(chan websites.Data)

	// Send a request to all websites registered
	for _, w := range m.websites {

		go func() {
			resultChan := w.Search(input.String())

			for {
				if res, ok := <-resultChan; ok {

					if movie, ok := res.(*websites.Movie); ok {
						c <- movie
					}
					log.Printf("continue!")
					continue
				}

				log.Printf("break!")
				break
			}

			close(c)
		}()
	}

	return c
}
