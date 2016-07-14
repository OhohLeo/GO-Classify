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

// Start the analysis of the import specified
func (c *Collection) Start() (chan *Item, error) {

	if c.items == nil {
		c.items = make(chan *Item)
	} else {
		return c.items, nil
	}

	// Start by analysing all imports
	for _, imported := range c.imports {
		c.startImport(imported)
	}

	return c.items, nil
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

// AddImport add new import
func (c *Collection) AddImport(name string, imported imports.Import) error {

	if _, ok := c.imports[name]; ok {
		return errors.New("import '" + name + "' already exists")
	}

	if c.imports == nil {
		c.imports = make(map[string]imports.Import)
	}

	c.imports[name] = imported

	if c.items != nil {
		c.startImport(imported)
	}

	return nil
}

// DeleteImport delete specified import
func (c *Collection) DeleteImport(name string) error {

	if imported, ok := c.imports[name]; ok {
		imported.Stop()
		delete(c.imports, name)
		return nil
	}

	return errors.New("import '" + name + "' not found")
}

// GetImports get the list of imports
func (c *Collection) GetImports() map[string]map[string]imports.Import {

	result := make(map[string]map[string]imports.Import)

	if c.imports == nil {
		return result
	}

	for name, imported := range c.imports {

		t := imported.GetType()

		if result[t] == nil {
			result[t] = make(map[string]imports.Import)
		}

		result[t][name] = imported
	}

	return result
}

// startImports launch the process of importation of specified import
func (c *Collection) startImport(imported imports.Import) error {

	// Get the import channel
	channel, err := imported.Start()
	if err != nil {
		return err
	}

	// Send all data imported to the collection
	go func() {
		for {
			if input, ok := <-channel; ok {
				c.items <- c.OnInput(input)
				continue
			}
			break
		}
	}()

	return nil
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

	// channel = make(chan websites.Data)

	// // Send a request to all websites registered
	// for _, w := range c.websites {

	// 	go func() {
	// 		resultChan := w.Search(input.String())

	// 		for {
	// 			if data, ok := <-resultChan; ok {
	// 				channel <- data
	// 				log.Printf("continue!")
	// 				continue
	// 			}

	// 			log.Printf("break!")
	// 			break
	// 		}

	// 		close(channel)
	// 	}()
	// }

	return
}
