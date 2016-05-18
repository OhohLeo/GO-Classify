package collections

import (
	"errors"
	"github.com/ohohleo/classify/imports"
	"github.com/ohohleo/classify/websites"
)

type Collection struct {
	imports map[string]imports.Import
	// exports map[exports.Export][]string
	websites map[string]websites.Website
}

// GetType returns the type of the collection (mandatory)
func (c *Collection) GetType() string {
	panic("Type should be specified!")
}

// AddImport add new import
func (c *Collection) AddImport(name string, imported imports.Import) {

	if c.imports == nil {
		c.imports = make(map[string]imports.Import)
	}

	c.imports[name] = imported
}

// Delete delete specified import
func (c *Collection) DeleteImport(name string) error {

	if _, ok := c.imports[name]; ok {
		delete(c.imports, name)
		return nil
	}

	return errors.New("no import name '" + name + "' found")
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
