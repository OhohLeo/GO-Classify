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
	buffer   *Buffer
	items    *Items
	config   *Config
	//exports map[exports.Export][]string
}

// SetName set the name of the collection
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

// SetName fix the name of the collection
func (c *Collection) SetConfig() {
	c.config = new(Config)
}

func (c *Collection) GetKeywords(item *Item) string {
	return "Star+Wars"
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

// WebResearch launch resarch through specified websites
func (c *Collection) WebResearch(item *Item) {

	// Get name
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

// OnInput handle new data to classify
func (c *Collection) OnInput(input imports.Data) *Item {

	// Create a new item
	item := NewItem()

	log.Printf("OnInput %s\n", item)

	// Add the import to the item
	item.AddImportData(input)

	// Store the import to the collection
	if c.buffer == nil {
		c.buffer = NewBuffer(2)
	}

	c.buffer.Add(item.GetId(), item)

	// // Launch research through web
	// if len(c.websites) > 0 {
	// 	c.WebResearch(item)
	// }

	// // TODO Research for best matching
	// item.Type = item.Websites["IMDB"][0].GetType()
	// item.IsMatching = 10.2
	// item.BestMatch = item.Websites["IMDB"][0]

	return item
}

func (c *Collection) GetBuffer() []*Item {
	return c.buffer.GetCurrentList()
}

func (c *Collection) CleanBuffer() {
	c.buffer = NewBuffer(c.config.BufferSize)
}
