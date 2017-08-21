package core

import (
	"encoding/json"
	"fmt"
	"github.com/ohohleo/classify/collections"
	"github.com/ohohleo/classify/database"
	"github.com/ohohleo/classify/imports"
	"github.com/ohohleo/classify/websites"
)

// Collection common methods
type Collection interface {
	Init(string, collections.Type) chan collections.Event
	SetName(string)
	GetName() string
	GetType() collections.Type
	ModifyConfig(string, string, []string) error
	ModifyConfigValue(string, string) error
	GetConfig() *collections.Config
	ResetBuffer()
	GetBuffer() []*collections.Item
	GetItems() []collections.Data
	Validate(string, *json.Decoder) error
	AddWebsite(website websites.Website)
	DeleteWebsite(name string) error
	OnInput(input imports.Data) *collections.Item
}

// Type of collections
var newCollections = []func() Collection{
	func() Collection { return new(collections.Movies) },
}

// Check collection names, returns the list of selected collections
func (c *Classify) GetCollectionsByNames(names []string) (map[string]Collection, error) {

	collections := make(map[string]Collection)

	for _, name := range names {

		collection, err := c.GetCollection(name)
		if err != nil {
			return nil, err
		}

		collections[name] = collection
	}

	return collections, nil
}

// Add a new collection
func (c *Classify) AddCollection(name string, collectionTypeStr string, toStore bool) (collection Collection, err error) {

	// Check that the name of the collection is unique
	if _, ok := c.collections[name]; ok {
		err = fmt.Errorf("collection '%s' already exists", name)
		return
	}

	// Check the collection type
	collectionType, ok := collections.TYPE_STR2IDX[collectionTypeStr]
	if ok == false || len(newCollections) <= int(collectionType) {
		err = fmt.Errorf("invalid collection type '%s'", collectionTypeStr)
		return
	}

	// Create the new collection
	new := newCollections[collectionType]()

	if c.collections == nil {
		c.collections = make(map[string]Collection)
	}

	// Store the collection
	c.collections[name] = new

	// Associate configuration
	eventsChannel := new.Init(name, collectionType)

	// Store the collection
	hDB, ok := new.(database.HandleDB)
	if toStore && ok {
		err = database.Insert(c.database, &collections.DB_DETAILS, hDB)
		if err != nil {
			fmt.Printf("=== ERROR %s \n", err.Error())
		}
	}

	go func() {

		for {
			event, ok := <-eventsChannel
			if ok {
				c.SendEvent("collection/"+name+"/"+event.Source,
					event.Status, event.Id, event.Item)
			}
		}
	}()

	return new, err
}

// Remove an existing collection
func (c *Classify) DeleteCollection(name string) (err error) {

	// Check that the name of the collection is unique
	if _, ok := c.collections[name]; ok == false {
		err = fmt.Errorf("collection '%s' not existing", name)
		return
	}

	delete(c.collections, name)
	return
}

// Return an existing collection
func (c *Classify) GetCollection(name string) (Collection, error) {

	collection, ok := c.collections[name]

	// Check that the name of the collection is unique
	if ok == false {
		return nil, fmt.Errorf("collection '%s' not existing", name)
	}

	return collection, nil
}

// Modify an existing collection
func (c *Classify) ModifyCollection(
	name string, newName string, newTypeStr string) (isModified bool, err error) {

	// Check that the name of the collection exists
	collection, ok := c.collections[name]
	if ok == false {
		err = fmt.Errorf("collection '%s' not existing", name)
		return
	}

	isModified = false

	if newTypeStr != "" {

		// Check the collection type
		newType, ok := collections.TYPE_STR2IDX[newTypeStr]
		if ok == false || len(newCollections) <= int(newType) {
			err = fmt.Errorf("invalid collection type '%s'", newTypeStr)
			return
		}

		if newType != collection.GetType() {

			// Check the collection type
			newCollection := newCollections[newType]

			delete(c.collections, name)
			collection = newCollection()
			c.collections[name] = newCollection()
			isModified = true
		}
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

// Returns the type of collection
func (c *Classify) GetCollectionTypes() []string {
	return collections.TYPE_IDX2STR
}
