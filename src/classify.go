package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/ohohleo/classify/collections"
	"github.com/ohohleo/classify/imports"
	"github.com/ohohleo/classify/websites"
)

type Classify struct {
	collections map[string]Collection
}

var classify *Classify

// Application startup
func Start() {

	classify = new(Classify)

	// TODO: Reload all collections saved

	log.SetLevel(log.DebugLevel)
	log.Info("Start Classify")
	ServerStart()
}

// Application stop
func (c *Classify) Stop() {
}

// Collection methods
type Collection interface {
	GetType() string
	Register(string, websites.Website)
	OnInput(imports.Data) chan websites.Data
}

var newCollections = map[string]func() Collection{
	"movies": func() Collection { return new(collections.Movies) },
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

// DeleteCollection remove an existing collection
func (c *Classify) DeleteCollection(name string) (err error) {

	// Check that the name of the collection is unique
	if _, ok := c.collections[name]; ok == false {
		err = fmt.Errorf("collection '%s' not existing", name)
		return
	}

	delete(c.collections, name)
	return
}

var collectionsList []string

// GetCollectionTypes returns the type of collections
func (c *Classify) GetCollectionTypes() []string {

	if collectionsList == nil {

		collectionsList = make([]string, len(newCollections))

		var id int = 0

		for name, _ := range newCollections {
			collectionsList[id] = name
			id++
		}
	}

	return collectionsList
}
