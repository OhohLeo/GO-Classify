package main

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/ohohleo/classify/collections"
	"github.com/ohohleo/classify/imports"
	"github.com/ohohleo/classify/imports/directory"
	"github.com/ohohleo/classify/websites"
	"github.com/ohohleo/classify/websites/IMDB"
)

// Collection methods
type Collection interface {
	Start() (chan *collections.Item, error)
	Stop()
	GetType() string
	AddImport(string, imports.Import) error
	DeleteImport(string) error
	GetImports() map[string]map[string]imports.Import
}

var newCollections = map[string]func() Collection{
	"movies": func() Collection { return new(collections.Movies) },
}

var newWebsites = map[string]websites.Website{
	"IMDB": IMDB.New(),
}

var newImports = map[string]func(json.RawMessage) (imports.Import, error){
	"directory": func(input json.RawMessage) (i imports.Import, err error) {
		var directory directory.Directory
		err = json.Unmarshal(input, &directory)
		if err == nil {
			i = &directory
		}
		return
	},
}

type Classify struct {
	collections map[string]Collection
}

var classify *Classify

// Application startup
func Start() *Classify {

	classify = new(Classify)

	// TODO: Reload all collections saved

	log.SetLevel(log.DebugLevel)
	log.Info("Start Classify")
	ServerStart()

	return classify
}

// Application stop
func (c *Classify) Stop() {
	ServerStop()
}

// AddCollection add a new collection
func (c *Classify) AddCollection(
	name string, collectionType string) (collection Collection, err error) {

	// Check that the name of the collection is unique
	if _, ok := c.collections[name]; ok {
		err = fmt.Errorf("collection '%s' already exists", name)
		return
	}

	// Get method to create new collection
	newCollection, ok := newCollections[collectionType]
	if ok == false {
		err = fmt.Errorf("invalid collection type '%s'", collectionType)
		return
	}

	// Create the new collection
	new := newCollection()

	if c.collections == nil {
		c.collections = make(map[string]Collection)
	}

	c.collections[name] = new

	return new, nil
}

// GetCollection return an existing collection
func (c *Classify) GetCollection(name string) (Collection, error) {

	collection, ok := c.collections[name]

	// Check that the name of the collection is unique
	if ok == false {
		return nil, fmt.Errorf("collection '%s' not existing", name)
	}

	return collection, nil
}

// ModifyCollection modify an existing collection
func (c *Classify) ModifyCollection(
	name string, newName string, newType string) (isModified bool, err error) {

	// Check that the name of the collection exists
	collection, ok := c.collections[name]
	if ok == false {
		err = fmt.Errorf("collection '%s' not existing", name)
		return
	}

	isModified = false

	if newType != "" && newType != collection.GetType() {

		// Check the collection type
		newCollection, ok := newCollections[newType]
		if ok == false {
			err = fmt.Errorf("invalid collection type '%s'", newType)
			return
		}

		delete(c.collections, name)
		collection = newCollection()
		c.collections[name] = newCollection()
		isModified = true
	}

	if newName != "" && newName != name {

		// Check that a collection called as newName doesn't exist
		_, ok := c.collections[newName]
		if ok {
			isModified = false
			err = fmt.Errorf("collection '%s' already existing", newName)
			return
		}

		delete(c.collections, name)
		c.collections[newName] = collection
		isModified = true
	}

	return
}

// DeleteCollection remove an existing exists
func (c *Classify) DeleteCollection(name string) (err error) {

	// Check that the name of the collection is unique
	if _, ok := c.collections[name]; ok == false {
		err = fmt.Errorf("collection '%s' not existing", name)
		return
	}

	delete(c.collections, name)
	return
}

var collectionTypes []string

// GetCollectionTypes returns the type of collection
func (c *Classify) GetCollectionTypes() []string {

	if collectionTypes == nil {

		collectionTypes = make([]string, len(newCollections))

		id := 0

		for name, _ := range newCollections {
			collectionTypes[id] = name
			id++
		}
	}

	return collectionTypes
}

var websitesList []string

// GetWebsites returns the list of websites available
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

type Mapping struct {
	Name   string          `json:"name"`
	Type   string          `json:"type"`
	Params json.RawMessage `json:"params"`
}

// CreateNewImport returns the import module depending on the parameter
// given
func (c *Classify) CreateNewImport(m Mapping) (name string, i imports.Import, err error) {

	if m.Name == "" {
		err = errors.New("Name field is mandatory!")
		return
	}

	name = m.Name

	if m.Type == "" {
		err = errors.New("Type field is mandatory!")
		return
	}

	createImport, ok := newImports[m.Type]
	if ok == false {
		err = errors.New("Import type '" + m.Type + "' not handled!")
		return
	}

	i, err = createImport(m.Params)
	if err != nil {
		return
	}

	return
}
