package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ohohleo/classify/database"
	"github.com/ohohleo/classify/exports"
	"github.com/ohohleo/classify/exports/file"
)

// Type of exports
var newExports = map[string]exports.Build{
	"file": file.ToBuild(),
}

type Export struct {
	Id          uint64 `json:"id"`
	Name        string `json:"name"`
	engine      exports.Export
	collections map[string]*Collection
}

func (i *Export) HasCollection(name string) (ok bool) {
	_, ok = i.collections[name]
	return
}

// Return true if export has a specified collection or no collections are specified
func (i *Export) HasCollections(collections map[string]*Collection) bool {

	if len(collections) > 0 {

		for name, _ := range collections {
			if i.HasCollection(name) {
				return true
			}
		}

		// No collection match
		return false
	}

	return true
}

func (e *Export) Store2DB(db *database.Database) error {

	// Check if db is enabled
	if db == nil {
		return nil
	}

	// Convert export to JSON
	paramsStr, err := json.Marshal(e.engine)
	if err != nil {
		return err
	}

	// Store the exports
	lastId, err := db.Insert("exports", &database.GenStruct{
		Name:   e.Name,
		Ref:    uint64(e.engine.GetRef()),
		Params: paramsStr,
	})
	if err != nil {
		return err
	}

	// Store the exports
	for _, collection := range e.collections {

		_, err := db.Insert("exports_mappings",
			map[string]interface{}{
				"exports_id":     lastId,
				"collections_id": collection.Id,
			})
		if err != nil {
			return err
		}
	}

	// Store current DB id
	e.Id = lastId

	return nil
}

func (e *Export) Unlink2DB(db *database.Database, collection *Collection) error {

	// Check if db is enabled
	if db == nil {
		return nil
	}

	return db.Delete("exports_mappings",
		map[string]interface{}{
			"exports_id":     e.Id,
			"collections_id": collection.Id},
		"exports_id = :exports_id AND collections_id = :collections_id")
}

func (e *Export) Delete2DB(db *database.Database) error {

	// Check if db is enabled
	if db == nil {
		return nil
	}

	return db.Delete("exports", &database.GenStruct{
		Id:  e.Id,
		Ref: uint64(e.engine.GetRef()),
	}, "id = :id AND ref = :ref")
}

// Check exports configuration
func (c *Classify) CheckExportsConfig(configuration map[string]json.RawMessage) (err error) {

	// For all export configuration
	for exportType, config := range configuration {

		// Select only generic type (ie ':type_name')
		if len(exportType) < 1 || exportType[0] != ':' {
			continue
		}

		// Check that the export type does exists
		buildExport, ok := newExports[exportType]
		if ok == false {
			err = errors.New("export type '" + exportType + "' not handled")
			return
		}

		// Check specified configuration
		err = buildExport.CheckConfig(config)
		if err != nil {
			return
		}
	}

	return nil
}

// Check exports ids and return the list of exports
func (c *Classify) GetExportsByNames(names []string) (exports map[string]*Export, err error) {

	exports = make(map[string]*Export)
	for _, name := range names {

		e, ok := c.exports[name]
		if ok == false {
			err = fmt.Errorf("import '%s' not found", name)
			return
		}

		exports[name] = e
	}

	return
}

// Add new export process
func (c *Classify) AddExport(name string, ref exports.Ref, params json.RawMessage, collections map[string]*Collection) (e *Export, err error) {

	// Nécessite l'existence d'au moins une collection
	if len(collections) < 1 {
		err = errors.New("required at least one existing collection")
		return
	}

	// Check that the ref exists
	buildExport, ok := newExports[ref.String()]
	if ok == false {
		err = errors.New("export ref '" + ref.String() + "' not handled")
		return
	}

	// TODO Get export configuration
	var config json.RawMessage

	// Get collections list
	idx := 0
	collectionNames := make([]string, len(collections))
	for name, _ := range collections {
		collectionNames[idx] = name
		idx++
	}

	// Create new export
	exportEngine, err := buildExport.Create(
		params, config, collectionNames)
	if err != nil {
		return
	}

	alreadyExists := false

	// Check if similar export already exists
	for _, e = range c.exports {

		// Returns similar export found
		if e.engine.GetRef() == ref && e.engine.Eq(exportEngine) {
			alreadyExists = true
			break
		}
	}

	// Otherwise create your export structure
	if alreadyExists == false {

		id := getRandomId()

		e = &Export{
			Id:          id,
			Name:        name,
			engine:      exportEngine,
			collections: collections,
		}

		if c.exports == nil {
			c.exports = make(map[string]*Export)
		}

		// Store the new export
		c.exports[name] = e

		return
	}

	e.collections = collections
	return
}

// Remove export from the list
func (c *Classify) DeleteExports(ids map[string]*Export, collections map[string]*Collection) (err error) {

	// At least one export id or one collection must be specified
	if len(ids) == 0 && len(collections) == 0 {
		err = errors.New("required export ids or collection names")
		return
	}

	// Stop all exports
	c.StopExports(ids, collections)

	// If no ids are specified : remove all export relative to the
	// same collection
	if len(ids) == 0 {
		ids = c.exports
	}

	for id, e := range ids {

		// Unlink the collection with the specified export
		for name, _ := range collections {
			delete(e.collections, name)
		}

		// If no collection are linked with specified export
		if len(e.collections) < 1 {

			if err = e.Delete2DB(c.database); err != nil {
				return
			}

			// Remove the export
			delete(c.exports, id)
		}
	}
	return
}

// Get the whole list of exports by Ref
func (c *Classify) GetExports(exportList map[string]*Export, collections map[string]*Collection) (res map[string]map[string]exports.Export, err error) {

	res = make(map[string]map[string]exports.Export)

	// If no exports are specified : get all
	if len(exportList) == 0 {
		exportList = c.exports
	}

	for name, e := range exportList {

		if e.HasCollections(collections) == false {
			continue
		}

		ref := e.engine.GetRef()

		if res[ref.String()] == nil {
			res[ref.String()] = make(map[string]exports.Export)
		}

		res[ref.String()][name] = e.engine
	}

	return
}

func (c *Classify) SendExportEvent(name string, status bool) {

	var statusStr string
	if status {
		statusStr = "start"
	} else {
		statusStr = "end"
	}

	c.SendEvent("export/status", statusStr, name, status)
}

// Force the exportation process on the collections specified
func (c *Classify) ForceExports(exportList map[string]*Export, collections map[string]*Collection) error {

	// If no exports are specified : get all
	if len(exportList) == 0 {
		exportList = c.exports
	}

	for name, e := range exportList {

		if e.HasCollections(collections) == false {
			continue
		}

		// For all collections
		for _, collection := range collections {

			// For all items
			for _, item := range collection.GetItems() {

				// Try to export them
				c.exports[name].engine.OnInput(item.Engine)
			}
		}

		// Send notification
		go c.SendExportEvent(name, false)
	}

	return nil
}

// Stop the exporting process
func (c *Classify) StopExports(exportList map[string]*Export, collections map[string]*Collection) error {

	// If no exports are specified : get all
	if len(exportList) == 0 {
		exportList = c.exports
	}

	for name, e := range exportList {

		if e.HasCollections(collections) == false {
			continue
		}

		c.exports[name].engine.Stop()

		// Send notification
		go c.SendExportEvent(name, false)
	}

	return nil
}

func (c *Classify) GetExportRefs() []string {
	return exports.REF_IDX2STR
}
