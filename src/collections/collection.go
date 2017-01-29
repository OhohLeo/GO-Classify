package collections

import (
	"errors"
	"github.com/ohohleo/classify/imports"
	"github.com/ohohleo/classify/websites"
	"log"
)

type Collection struct {
	name     string
	websites map[string]websites.Website
	items    map[string]*Item
	//exports map[exports.Export][]string
}

// SetName fix the name of the collection
func (c *Collection) SetName(name string) {
	c.name = name
}

// GetName returns the name of the collection
func (c *Collection) GetName() string {
	return c.name
}

// GetType returns the type of the collection (mandatory)
func (c *Collection) GetType() string {
	panic("collection type should be specified!")
}

func (c *Collection) GetKeywords(item *Item) string {
	return item.GetKeywords()
}

// AddWebsite add new website
func (c *Collection) AddWebsite(website websites.Website) {

	if c.websites == nil {
		c.websites = make(map[string]websites.Website)
	}

	c.websites[website.GetName()] = website
}

// DeleteWebsite delete specified website
func (c *Collection) DeleteWebsite(name string) error {

	if _, ok := c.websites[name]; ok {

		delete(c.websites, name)
		return nil
	}

	return errors.New("no website name '" + name + "' found")
}

// OnInput handle new data to classify
func (c *Collection) OnInput(input imports.Data) *Item {

	// Create a new item
	item := NewItem()

	log.Printf("OnInput %s\n", item)

	// Add the import to the item
	item.AddImportData(input)

	// Store the import to the collection
	if c.items == nil {
		c.items = make(map[string]*Item)
	}

	c.items[item.GetId()] = item

	// Launch research through web
	if len(c.websites) > 0 {
		c.WebResearch(item)
	}

	// TODO Research for best matching
	item.Type = item.Websites["IMDB"][0].GetType()
	item.IsMatching = 10.2
	item.BestMatch = item.Websites["IMDB"][0]

	return item
}

// WebResearch launch resarch through specified websites
func (c *Collection) WebResearch(item *Item) {

	keywords := c.GetKeywords(item)

	// For all specified websites
	for _, website := range c.websites {

		// Launch the research
		channel := website.Search(keywords)

		for {
			data, ok := <-channel
			if ok {
				//log.Printf("WEB DATA %+v\n", data)
				item.AddWebsiteData(website.GetName(), data)
			}

			break
		}
	}

}
