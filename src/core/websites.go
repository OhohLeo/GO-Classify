package core

import (
	"errors"
	"github.com/ohohleo/classify/websites"
	"github.com/ohohleo/classify/websites/IMDB"
)

var newWebsites = map[string]websites.Website{
	"IMDB": IMDB.New(),
}

// AddWebsite add new website
func (c *Classify) AddWebsite(name string) (website websites.Website, err error) {

	var ok bool

	// Check if the website already exits
	if c.websites != nil {
		if website, ok = c.websites[name]; ok {
			return
		}
	} else {
		c.websites = make(map[string]websites.Website)
	}

	website, ok = newWebsites[name]
	if ok == false {
		err = errors.New("no existing website '" + name + "'")
		return
	}

	c.websites[name] = website
	return
}

// DeleteWebsite delete specified website
func (c *Classify) DeleteWebsite(name string) error {

	if _, ok := c.websites[name]; ok {
		delete(c.websites, name)
		return nil
	}

	return errors.New("no website name '" + name + "' found")
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
