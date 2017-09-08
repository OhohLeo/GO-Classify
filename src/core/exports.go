package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ohohleo/classify/database"
	"github.com/ohohleo/classify/exports"
	"github.com/ohohleo/classify/exports/file"
	"strconv"
)

type BuildExport struct {
	CheckConfig func(config map[string][]string) error
	Create      func(json.RawMessage, map[string][]string, []string) (exports.Export, error)
}

// Type of exports
var newExports = map[string]exports.BuildExport{
	"file": file.ToBuild(),
}

type Export struct {
	Id          uint64 `json:"id"`
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
func (c *Classify) CheckExportsConfig(configuration map[string]map[string][]string) (err error) {

	// For all export configuration
	for exportRef, config := range configuration {

		// Check that the export ref does exists
		buildExport, ok := newExports[exportRef]
		if ok == false {
			err = errors.New("export ref '" + exportRef + "' not handled")
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
func (c *Classify) GetExportsByIds(ids []uint64) (exports map[uint64]*Export, err error) {

	exports = make(map[uint64]*Export)

	for _, id := range ids {
		i, ok := c.exports[id]
		if ok == false {
			err = fmt.Errorf("export '%s' not existing", id)
			return
		}

		exports[id] = i
	}

	return
}

// Add new export process
func (c *Classify) AddExport(ref exports.Ref, params json.RawMessage, collections map[string]*Collection) (e *Export, err error) {

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

	// Get export configuration
	var config map[string][]string
	if c.config != nil {
		config, _ = c.config.Exports[ref.String()]
	}

	// Get collections list
	idx := 0
	collectionNames := make([]string, len(collections))
	for name, _ := range collections {
		collectionNames[idx] = name
		idx++
	}

	// Create new export
	exportEngine, err := buildExport.Create(params, config, collectionNames)
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
			engine:      exportEngine,
			collections: collections,
		}

		fmt.Printf("NEW EXPORT %+v\n", e)

		if c.exports == nil {
			c.exports = make(map[uint64]*Export)
		}

		// Store the new export
		c.exports[id] = e

		return
	}

	e.collections = collections
	return
}

// Remove export from the list
func (c *Classify) DeleteExports(ids map[uint64]*Export, collections map[string]*Collection) (err error) {

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
func (c *Classify) GetExports(ids map[uint64]*Export, collections map[string]*Collection) (res map[string]map[uint64]exports.Export, err error) {

	res = make(map[string]map[uint64]exports.Export)

	// If no ids are specified : get all
	if len(ids) == 0 {
		ids = c.exports
	}

	for name, e := range ids {

		if e.HasCollections(collections) == false {
			continue
		}

		ref := e.engine.GetRef()

		if res[ref.String()] == nil {
			res[ref.String()] = make(map[uint64]exports.Export)
		}

		res[ref.String()][name] = e.engine
	}

	return
}

func (c *Classify) SendExportEvent(id uint64, status bool) {

	var statusStr string
	if status {
		statusStr = "start"
	} else {
		statusStr = "end"
	}

	c.SendEvent("export/status", statusStr, strconv.Itoa(int(id)), status)
}

// Force the exportation process on the collections specified
func (c *Classify) ForceExports(ids map[uint64]*Export, collections map[string]*Collection) error {

	// If no ids are specified : get all
	if len(ids) == 0 {
		ids = c.exports
	}

	for id, e := range ids {

		if e.HasCollections(collections) == false {
			continue
		}

		// For all collections
		for _, collection := range collections {

			// For all items
			for _, item := range collection.GetItems() {

				// Try to export them
				c.exports[id].engine.OnInput(item.Engine)
			}
		}

		// Send notification
		go c.SendExportEvent(id, false)
	}

	return nil
}

// Stop the exporting process
func (c *Classify) StopExports(ids map[uint64]*Export, collections map[string]*Collection) error {

	// If no ids are specified : get all
	if len(ids) == 0 {
		ids = c.exports
	}

	for id, e := range ids {

		if e.HasCollections(collections) == false {
			continue
		}

		c.exports[id].engine.Stop()

		// Send notification
		go c.SendExportEvent(id, false)
	}

	return nil
}

func (c *Classify) GetExportRefs() []string {
	return exports.REF_IDX2STR
}
