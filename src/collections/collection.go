package collections

import (
	"errors"
	"github.com/ohohleo/classify/imports"
	"github.com/ohohleo/classify/websites"
	"log"
)

type CollectionData struct {
	Import []imports.Data
}

type Collection struct {
	imports map[string]imports.Import
	// exports map[exports.Export][]string
	websites map[string]websites.Website
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

	return nil
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

// LaunchImport starts the import specified
func (c *Collection) LaunchImport(name string) (chan imports.Data, error) {

	res := make(chan imports.Data)

	// Get the import item
	imported, ok := c.imports[name]
	if ok == false {
		return nil, errors.New("import '" + name + "' not found")
	}

	// Get the import channel
	channel, err := imported.Launch()
	if err != nil {
		return nil, err
	}

	// Send all data imported to the collection
	go func() {
		for {
			if input, ok := <-channel; ok {
				res <- input
				c.OnInput(input)
				log.Printf("COLLECTION %+v\n", input.GetType())
				continue
			}
			break
		}
	}()

	return res, nil
}

// Delete delete specified import
func (c *Collection) DeleteImport(name string) error {

	if _, ok := c.imports[name]; ok {
		delete(c.imports, name)
		return nil
	}

	return errors.New("import '" + name + "' not found")
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
func (c *Collection) OnInput(input imports.Data) chan websites.Data {

	channel := make(chan websites.Data)

	// Send a request to all websites registered
	for _, w := range c.websites {

		go func() {
			resultChan := w.Search(input.String())

			for {
				if data, ok := <-resultChan; ok {
					channel <- data
					log.Printf("continue!")
					continue
				}

				log.Printf("break!")
				break
			}

			close(channel)
		}()
	}

	return channel
}
