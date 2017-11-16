package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ohohleo/classify/collections"
	// "github.com/ohohleo/classify/websites"
	"log"
)

// Type of collections
var newCollections = map[string]collections.Build{
	"movies": collections.BuildMovies(),
	"simple": collections.BuildSimple(),
}

// Check imports configuration
func (c *Classify) CheckCollectionsConfig(configuration map[string]json.RawMessage) (err error) {

	// For all collection configuration
	for collectionType, config := range configuration {

		// Select only generic type (ie ':type_name')
		if len(collectionType) < 1 || collectionType[0] != ':' {
			continue
		}

		// Check that the collection type does exists
		buildCollection, ok := newCollections[collectionType]
		if ok == false {
			err = errors.New("collection type '" + collectionType + "' not handled")
			return
		}

		// Check specified configuration
		err = buildCollection.CheckConfig(config)
		if err != nil {
			return
		}
	}

	return nil
}

// Check collection names, returns the list of selected collections
func (c *Classify) GetCollectionsByNames(names []string) (map[string]*Collection, error) {

	collections := make(map[string]*Collection)

	for _, name := range names {

		collection, err := c.GetCollection(name)
		if err != nil {
			return nil, err
		}

		collections[name] = collection
	}

	return collections, nil
}

func (c *Classify) AddCollection(name string, ref collections.Ref, config json.RawMessage, params json.RawMessage) (collection *Collection, err error) {

	// var website websites.Website

	// // Check if the websites does exists
	// websites := make([]websites.Website, 0)
	// for _, name := range webNames {

	// 	website, err = c.AddWebsite(name)
	// 	if err != nil {
	// 		return
	// 	}

	// 	// Add new website
	// 	websites = append(websites, website)
	// }

	// Add stored collection
	collection, err = c.addCollection(name, ref, config, params)
	if err != nil {
		return
	}

	log.Printf("Add collection '%s' as '%s'\n", name, ref.String())

	// // Add websites to the collection created
	// for _, website := range websites {
	// 	collection.AddWebsite(website)
	// }

	return
}

// Add a new collection
func (c *Classify) addCollection(name string, ref collections.Ref, config json.RawMessage, params json.RawMessage) (collection *Collection, err error) {

	// Check that the name of the collection is unique
	if _, ok := c.collections[name]; ok {
		err = fmt.Errorf("collection '%s' already exists", name)
		return
	}

	buildCollection, ok := newCollections[ref.String()]
	if ok == false {
		err = errors.New("collection type '" + ref.String() + "' not handled")
		return
	}

	// Create the new collection
	collectionEngine, err := buildCollection.Create(params, nil)
	if err != nil {
		return
	}

	eventsChannel := make(chan CollectionEvent)

	collection = &Collection{
		Name:   name,
		items:  NewItems(),
		config: NewCollectionConfig(),
		engine: collectionEngine,
	}

	// Store configuration received
	if config != nil {
		err = json.Unmarshal(config, collection.config)
		if err != nil {
			return
		}
	}

	if c.collections == nil {
		c.collections = make(map[string]*Collection)
	}

	// Store the collection
	c.collections[name] = collection

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
func (c *Classify) GetCollection(name string) (*Collection, error) {

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

		if newRef != collection.engine.GetRef() {

			// // Check the collection type
			// newCollection := newCollections[newRef]

			// delete(c.collections, name)
			// collection = newCollection()
			// c.collections[name] = newCollection()
			// isModified = true
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

		// delete(c.collections, name)
		// c.collections[newName] = collection
		// isModified = true
	}

	return
}

// Returns the type of collection
func (c *Classify) GetCollectionRefs() []string {
	return collections.REF_IDX2STR
}
