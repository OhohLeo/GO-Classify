package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ohohleo/classify/database"
	"github.com/ohohleo/classify/imports"
	"github.com/ohohleo/classify/imports/directory"
	"github.com/ohohleo/classify/imports/email"
	"strconv"
)

// Type of imports
var newImports = map[string]imports.BuildImport{
	"directory": directory.ToBuild(),
	"email":     email.ToBuild(),
}

type Import struct {
	Id          string `json:"id"`
	engine      imports.Import
	collections map[string]Collection
}

func (i *Import) HasCollection(name string) (ok bool) {
	_, ok = i.collections[name]
	return
}

// Return true if import has a specified collection or no collections are specified
func (i *Import) HasCollections(collections map[string]Collection) bool {

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

func (i *Import) Store2DB(db *database.Database) error {

	// Check if db is enabled
	if db == nil {
		return nil
	}

	// Convert import to JSON
	paramsStr, err := json.Marshal(i.engine)
	if err != nil {
		return err
	}

	// Store the imports
	lastId, err := db.Insert("imports", &database.GenStruct{
		Ref:    uint64(i.engine.GetRef()),
		Params: paramsStr,
	})
	if err != nil {
		return err
	}

	// Store the imports
	for _, collection := range i.collections {

		_, err := db.Insert("imports_mappings",
			map[string]interface{}{
				"imports_id":     lastId,
				"collections_id": collection.GetId(),
			})
		if err != nil {
			return err
		}
	}

	// Store current DB id
	i.Id = strconv.Itoa(int(lastId))

	return nil
}

func (i *Import) Unlink2DB(db *database.Database, collection Collection) error {

	// Check if db is enabled
	if db == nil {
		return nil
	}

	// Get Id
	id, err := strconv.Atoi(i.Id)
	if err != nil {
		return err
	}

	return db.Delete("imports_mappings",
		map[string]interface{}{
			"imports_id":     id,
			"collections_id": collection.GetId()}, "id = :id")
}

func (i *Import) Delete2DB(db *database.Database) error {

	// Check if db is enabled
	if db == nil {
		return nil
	}

	id, err := strconv.Atoi(i.Id)
	if err != nil {
		return err
	}

	return db.Delete("imports", &database.GenStruct{
		Id:  uint64(id),
		Ref: uint64(i.engine.GetRef()),
	}, "id = :id AND ref = :ref")
}

// Check imports configuration
func (c *Classify) CheckImportsConfig(configuration map[string]map[string][]string) (err error) {

	// For all import configuration
	for importType, config := range configuration {

		// Check that the import type does exists
		buildImport, ok := newImports[importType]
		if ok == false {
			err = errors.New("import type '" + importType + "' not handled")
			return
		}

		// Check specified configuration
		err = buildImport.CheckConfig(config)
		if err != nil {
			return
		}
	}

	return nil
}

// Check imports ids and return the list of imports
func (c *Classify) GetImportsByIds(ids []string) (imports map[string]Import, err error) {

	imports = make(map[string]Import)

	for _, id := range ids {
		i, ok := c.imports[id]
		if ok == false {
			err = fmt.Errorf("import '%s' not existing", id)
			return
		}

		imports[id] = i
	}

	return
}

// Add new import process
func (c *Classify) AddImport(importRef string, params json.RawMessage, collections map[string]Collection) (i Import, err error) {

	// Nécessite l'existence d'au moins une collection
	if len(collections) < 1 {
		err = errors.New("required at least one existing collection")
		return
	}

	// Field required
	if importRef == "" {
		err = errors.New("type field is mandatory")
		return
	}

	// Check that the type exists
	buildImport, ok := newImports[importRef]
	if ok == false {
		err = errors.New("import type '" + importRef + "' not handled")
		return
	}

	// Get import configuration
	// config, _ := c.config.Imports[importRef]
	var config map[string][]string

	// Get collections list
	idx := 0
	collectionNames := make([]string, len(collections))
	for name, _ := range collections {
		collectionNames[idx] = name
		idx++
	}

	// Create new import
	importEngine, err := buildImport.Create(params, config, collectionNames)
	if err != nil {
		return
	}

	alreadyExists := false

	// Check if similar import already exists
	for _, i = range c.imports {

		// Returns similar import found
		if i.engine.GetRef().String() == importRef && i.engine.Eq(importEngine) {
			alreadyExists = true
			break
		}
	}

	// Otherwise create your import structure
	if alreadyExists == false {

		id := getRandomName()
		i = Import{
			Id:          id,
			engine:      importEngine,
			collections: collections,
		}

		if c.imports == nil {
			c.imports = make(map[string]Import)
		}

		// Store the new import
		c.imports[id] = i

		return
	}

	i.collections = collections
	return
}

// Remove import from the list
func (c *Classify) DeleteImports(ids map[string]Import, collections map[string]Collection) (err error) {

	// At least one import id or one collection must be specified
	if len(ids) == 0 && len(collections) == 0 {
		err = errors.New("required import ids or collection names")
		return
	}

	// Stop all imports
	c.StopImports(ids, collections)

	// If no ids are specified : remove all import relative to the
	// same collection
	if len(ids) == 0 {
		ids = c.imports
	}

	for id, i := range ids {

		// Unlink the collection with the specified import
		for name, collection := range collections {

			if err = i.Unlink2DB(c.database, collection); err != nil {
				return
			}

			delete(i.collections, name)
		}

		// If no collection are linked with specified import
		if len(i.collections) < 1 {

			if err = i.Delete2DB(c.database); err != nil {
				return
			}

			// Remove the import
			delete(c.imports, id)
		}
	}
	return
}

// Get the whole list of imports by Type
func (c *Classify) GetImports(ids map[string]Import, collections map[string]Collection) (res map[string]map[string]imports.Import, err error) {

	res = make(map[string]map[string]imports.Import)

	// If no ids are specified : get all
	if len(ids) == 0 {
		ids = c.imports
	}

	for name, i := range ids {

		if i.HasCollections(collections) == false {
			continue
		}

		t := i.engine.GetRef()

		if res[t.String()] == nil {
			res[t.String()] = make(map[string]imports.Import)
		}

		res[t.String()][name] = i.engine
	}

	return
}

func (c *Classify) SendImportEvent(id string, status bool) {

	var statusStr string
	if status {
		statusStr = "start"
	} else {
		statusStr = "end"
	}

	c.SendEvent("import/status", statusStr, id, status)
}

// Launch the process of importation of specified imports
func (c *Classify) StartImports(ids map[string]Import, collections map[string]Collection) error {

	// If no ids are specified : get all
	if len(ids) == 0 {
		ids = c.imports
	}

	// Get the import channel
	for id, i := range ids {

		if i.HasCollections(collections) == false {
			continue
		}

		channel, err := i.engine.Start()
		if err != nil {
			return err
		}

		// Send all data imported to the collections
		go func() {

			// Send notification to start analysis
			c.SendImportEvent(id, true)

			for {
				if input, ok := <-channel; ok {

					// For each collections linked with the importation
					for _, collection := range i.collections {
						collection.OnInput(input)
					}

					continue
				}

				break
			}

			// Send notification to stop analysis
			c.SendImportEvent(id, false)
		}()
	}

	return nil
}

// Stop the importing process
func (c *Classify) StopImports(ids map[string]Import, collections map[string]Collection) error {

	// If no ids are specified : get all
	if len(ids) == 0 {
		ids = c.imports
	}

	for id, i := range ids {

		if i.HasCollections(collections) == false {
			continue
		}

		c.imports[id].engine.Stop()

		// Send notification
		go c.SendImportEvent(id, false)
	}

	return nil
}
