package core

import (
	"github.com/ohohleo/classify/websites"
	"github.com/ohohleo/classify/websites/IMDB"
)

var newWebsites = map[string]websites.Website{
	"IMDB": IMDB.New(),
}

var websitesList []string

// Returns the list of websites available
func (c *Classify) GetWebsites() []string {

	if websitesList == nil {

		websitesList = make([]string, len(newWebsites))

		id := 0

		for name, _ := range newWebsites {
			websitesList[id] = name
			id++
		}
	}

	return websitesList
}
