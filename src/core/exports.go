package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ohohleo/classify/exports"
	"strconv"
)

type BuildExport struct {
	CheckConfig func(config map[string][]string) error
	Create      func(json.RawMessage, map[string][]string, []string) (exports.Export, error)
}

// Type of exports
var newExports = map[string]BuildExport{}

type Export struct {
	Id          uint64 `json:"id"`
	engine      exports.Export
	collections map[string]Collection
}

func (i *Export) HasCollection(name string) (ok bool) {
	_, ok = i.collections[name]
	return
}

// Return true if export has a specified collection or no collections are specified
func (i *Export) HasCollections(collections map[string]Collection) bool {

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

// Check exports configuration
func (c *Classify) CheckExportsConfig(configuration map[string]map[string][]string) (err error) {

	// For all export configuration
	for exportType, config := range configuration {

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
func (c *Classify) GetExportsByIds(ids []uint64) (exports map[uint64]Export, err error) {

	exports = make(map[uint64]Export)

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
func (c *Classify) AddExport(exportType string, params json.RawMessage, collections map[string]Collection) (i Export, err error) {

	// Nécessite l'existence d'au moins une collection
	if len(collections) < 1 {
		err = errors.New("required at least one existing collection")
		return
	}

	// Field required
	if exportType == "" {
		err = errors.New("type field is mandatory")
		return
	}

	// Check that the type exists
	buildExport, ok := newExports[exportType]
	if ok == false {
		err = errors.New("export type '" + exportType + "' not handled")
		return
	}

	// Get export configuration
	config, _ := c.config.Exports[exportType]

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
	for _, i = range c.exports {

		// Returns similar export found
		if i.engine.GetType() == exportType && i.engine.Eq(exportEngine) {
			alreadyExists = true
			break
		}
	}

	// Otherwise create your export structure
	if alreadyExists == false {

		id := getRandomId()
		i = Export{
			Id:          id,
			engine:      exportEngine,
			collections: collections,
		}

		if c.exports == nil {
			c.exports = make(map[uint64]Export)
		}

		// Store the new export
		c.exports[id] = i

		return
	}

	i.collections = collections
	return
}

// Remove export from the list
func (c *Classify) DeleteExports(ids map[uint64]Export, collections map[string]Collection) (err error) {

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

	for id, i := range ids {

		// Unlink the collection with the specified export
		for name, _ := range collections {
			delete(i.collections, name)
		}

		// If no collection are linked with specified export
		if len(i.collections) < 1 {

			// Remove the export
			delete(c.exports, id)
		}
	}
	return
}

// Get the whole list of exports by Type
func (c *Classify) GetExports(ids map[uint64]Export, collections map[string]Collection) (res map[string]map[uint64]exports.Export, err error) {

	res = make(map[string]map[uint64]exports.Export)

	// If no ids are specified : get all
	if len(ids) == 0 {
		ids = c.exports
	}

	for name, i := range ids {

		if i.HasCollections(collections) == false {
			continue
		}

		t := i.engine.GetType()

		if res[t] == nil {
			res[t] = make(map[uint64]exports.Export)
		}

		res[t][name] = i.engine
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

// Stop the exporting process
func (c *Classify) StopExports(ids map[uint64]Export, collections map[string]Collection) error {

	// If no ids are specified : get all
	if len(ids) == 0 {
		ids = c.exports
	}

	for id, i := range ids {

		if i.HasCollections(collections) == false {
			continue
		}

		c.exports[id].engine.Stop()

		// Send notification
		go c.SendExportEvent(id, false)
	}

	return nil
}
