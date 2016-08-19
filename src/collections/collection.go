package collections

import (
	"errors"
	"github.com/ohohleo/classify/imports"
	"github.com/ohohleo/classify/websites"
	//	"log"
)

type Collection struct {
	name     string
	imports  map[string]imports.Import
	websites map[string]websites.Website
	// exports map[exports.Export][]string

	importsToItem map[string]*Item
	items         chan *Item
}

func (c *Collection) initItems() chan *Item {

	if c.items == nil {
		c.items = make(chan *Item)
	}

	return c.items
}

// Start the analysis of the import specified
func (c *Collection) Start() (chan *Item, error) {

	items := c.initItems()

	// Start by analysing all imports
	for _, imported := range c.imports {
		c.startImport(imported)
	}

	return items, nil
}

// Stop analysis of the collection
func (c *Collection) Stop() {

	// Stop imported process
	for _, imported := range c.imports {
		imported.Stop()
	}
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

// AddExport add new export process
func (c *Collection) AddExport(name string) {
}

// DeleteExport delete specified export
func (c *Collection) DeleteExport(name string) {
}

// AddWebsite add new website
func (c *Collection) AddWebsite(name string, website websites.Website) {

	if c.websites == nil {
		c.websites = make(map[string]websites.Website)
	}

	c.websites[name] = website
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
func (c *Collection) OnInput(input imports.Data) (item *Item) {

	// Check if a similar input doesn't exist yet
	inputKey := input.GetUniqKey()

	if _, ok := c.importsToItem[inputKey]; ok {
		// No need to add similar input
		return
	}

	// Create a new item
	item = new(Item)

	// Add the input to the item
	item.AddImportData(input)

	// Store the inputs to the collection
	if c.importsToItem == nil {
		c.importsToItem = make(map[string]*Item)
	}

	c.importsToItem[inputKey] = item

	return
}
