package core

import (
	"encoding/json"
	"fmt"
	"github.com/ohohleo/classify/collections"
	"github.com/ohohleo/classify/database"
	"github.com/ohohleo/classify/imports"
	"github.com/ohohleo/classify/websites"
)

// type Collection struct {
// 	Id     uint64 `json:"id"`
// 	engine collections.Collection
// }

// Collection common methods
type Collection interface {
	Init(string, collections.Ref) chan collections.Event
	SetId(uint64)
	GetId() uint64
	SetName(string)
	GetName() string
	GetRef() collections.Ref
	ModifyConfig(string, string, []string) error
	ModifyConfigValue(string, string) error
	GetConfig() *collections.Config
	ResetBuffer()
	GetBuffer() []*collections.BufferItem
	Validate(string, *json.Decoder) error
	AddWebsite(website websites.Website)
	DeleteWebsite(name string) error
	OnInput(input imports.Data) *collections.BufferItem
	Store2DB(db *database.Database) error
	Delete2DB(db *database.Database) error
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

func (c *Classify) AddCollection(name string, ref collections.Ref, webNames []string) (collection Collection, err error) {

	var website websites.Website

	// Check if the websites does exists
	websites := make([]websites.Website, 0)
	for _, name := range webNames {

		website, err = c.AddWebsite(name)
		if err != nil {
			return
		}

		// Add new website
		websites = append(websites, website)
	}

	// Add stored collection
	collection, err = c.addCollection(name, ref)
	if err != nil {
		return
	}

	// Add websites to the collection created
	for _, website := range websites {
		collection.AddWebsite(website)
	}

	return
}

// Add a new collection
func (c *Classify) addCollection(name string, ref collections.Ref) (collection Collection, err error) {

	// Check that the name of the collection is unique
	if _, ok := c.collections[name]; ok {
		err = fmt.Errorf("collection '%s' already exists", name)
		return
	}

	// Create the new collection
	collection = newCollections[ref]()

	if c.collections == nil {
		c.collections = make(map[string]Collection)
	}

	// Store the collection
	c.collections[name] = collection

	// Associate configuration
	eventsChannel := collection.Init(name, ref)

	go func() {

		for {
			event, ok := <-eventsChannel
			if ok {
				c.SendEvent("collection/"+name+"/"+event.Source,
					event.Status, event.Id, event.Item)
			}
		}
	}()

	return
}

// Remove an existing collection
func (c *Classify) DeleteCollection(name string) (err error) {

	collection, ok := c.collections[name]

	// Check that the name of the collection is unique
	if ok == false {
		err = fmt.Errorf("collection '%s' not existing", name)
		return
	}

	// Store collection if enable
	if c.database != nil {
		err = collection.Delete2DB(c.database)
		if err != nil {
			err = fmt.Errorf("delete collection '%s' from DB failed: %s",
				name, err.Error())
			return
		}
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
	name string, newName string, newRefStr string) (isModified bool, err error) {

	// Check that the name of the collection exists
	collection, ok := c.collections[name]
	if ok == false {
		err = fmt.Errorf("collection '%s' not existing", name)
		return
	}

	isModified = false

	if newRefStr != "" {

		// Check the collection type
		newRef, ok := collections.REF_STR2IDX[newRefStr]
		if ok == false || len(newCollections) <= int(newRef) {
			err = fmt.Errorf("invalid collection type '%s'", newRefStr)
			return
		}

		if newRef != collection.GetRef() {

			// Check the collection type
			newCollection := newCollections[newRef]

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
func (c *Classify) GetCollectionRefs() []string {
	return collections.REF_IDX2STR
}
