package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/ohohleo/classify/data"
	"github.com/ohohleo/classify/database"
	"github.com/ohohleo/classify/imports"
	"github.com/ohohleo/classify/imports/directory"
	"github.com/ohohleo/classify/imports/imap"
)

// Type of imports
var newImports = map[string]imports.Build{
	"directory": directory.ToBuild(),
	"imap":      imap.ToBuild(),
}

type Import struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`

	engine imports.Import

	collectionImports map[string]*CollectionImport
}

func (i *Import) GetDatas() map[string]interface{} {

	res := make(map[string]interface{})

	for _, data := range i.engine.GetDataList() {
		res[data.GetRef().String()] = data
	}

	return res
}

func (i *Import) HasCollection(collectionName string) (ok bool) {
	_, ok = i.collectionImports[collectionName]
	return
}

// Return true if import has a specified collection or no collections are specified
func (i *Import) HasCollections(collections map[string]*Collection) bool {

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

func (i *Import) GetConfig(collectionName string) (config *ImportExportConfig, err error) {

	if collection, ok := i.collectionImports[collectionName]; ok {
		config = collection.Config
	} else {
		err = fmt.Errorf("no config found for collection '%s'", collectionName)
	}
	return
}

func (i *Import) SetConfig(collectionName string, newConfig *ImportExportConfig) (err error) {

	var actualConfig *ImportExportConfig
	actualConfig, err = i.GetConfig(collectionName)
	if err != nil {
		return
	}

	// TODO : handle enable/disable import
	if actualConfig.General.Enabled != newConfig.General.Enabled {
	}

	if newConfig.Tweak != nil {
		actualConfig.SetTweak(newConfig.Tweak)
	}

	return
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
		Name:   i.Name,
		Ref:    uint64(i.engine.GetRef()),
		Params: paramsStr,
	})
	if err != nil {
		return err
	}

	// Store the imports
	for _, collectionImport := range i.collectionImports {

		_, err := db.Insert("imports_mappings",
			map[string]interface{}{
				"imports_id":     lastId,
				"collections_id": collectionImport.Collection.Id,
			})
		if err != nil {
			return err
		}
	}

	// Store current DB id
	i.Id = lastId

	return nil
}

func (i *Import) Unlink2DB(db *database.Database, collection *Collection) error {

	// Check if db is enabled
	if db == nil {
		return nil
	}

	return db.Delete("imports_mappings",
		map[string]interface{}{
			"imports_id":     i.Id,
			"collections_id": collection.Id},
		"imports_id = :imports_id AND collections_id = :collections_id")
}

func (i *Import) Delete2DB(db *database.Database) error {

	// Check if db is enabled
	if db == nil {
		return nil
	}

	return db.Delete("imports", &database.GenStruct{
		Id:  i.Id,
		Ref: uint64(i.engine.GetRef()),
	}, "id = :id AND ref = :ref")
}

// Check imports configuration
func (c *Classify) CheckImportsConfig(configuration map[string]json.RawMessage) (err error) {

	// For all import configuration
	for importType, config := range configuration {

		// Select only generic type (ie ':type_name')
		if len(importType) < 1 || importType[0] != ':' {
			continue
		}

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

// Check import name and return the imports
func (c *Classify) GetImportByName(name string) (i *Import, err error) {

	var ok bool
	i, ok = c.imports[name]
	if ok == false {
		err = fmt.Errorf("import '%s' not found", name)
		return
	}

	return
}

// Check imports names and return the list of imports
func (c *Classify) GetImportsByNames(names []string) (imports map[string]*Import, err error) {

	imports = make(map[string]*Import)
	for _, name := range names {
		imports[name], err = c.GetImportByName(name)
		if err != nil {
			return
		}
	}

	return
}

// Add new import process
func (c *Classify) AddImport(name string, ref imports.Ref, inParams json.RawMessage, collections map[string]*Collection) (i *Import, outParams interface{}, err error) {

	// Nécessite l'existence d'au moins une collection
	if len(collections) < 1 {
		err = errors.New("required at least one existing collection")
		return
	}

	// Check that the type exists
	buildImport, ok := newImports[ref.String()]
	if ok == false {
		err = errors.New("import type '" + ref.String() + "' not handled")
		return
	}

	// TODO Get import configuration
	var config json.RawMessage

	// Get collections list
	idx := 0
	collectionNames := make([]string, len(collections))
	for name, _ := range collections {
		collectionNames[idx] = name
		idx++
	}

	// Create new import
	importEngine, moreParams, err := buildImport.Create(
		inParams, config, collectionNames)
	if err != nil {
		outParams = moreParams
		return
	}

	alreadyExists := false

	// Check if similar import already exists
	for _, i = range c.imports {

		// Returns similar import found
		if i.engine.GetRef() == ref && i.engine.Eq(importEngine) {
			alreadyExists = true
			break
		}
	}

	// Otherwise create your import structure
	if alreadyExists == false {

		id := getRandomId()

		i = &Import{
			Id:     id,
			Name:   name,
			engine: importEngine,
		}

		if c.imports == nil {
			c.imports = make(map[string]*Import)
		}

		// Store the new import
		c.imports[name] = i
	}

	// Set collection list to the import
	i.collectionImports = make(map[string]*CollectionImport)

	// Add import to the collection
	for _, collection := range collections {

		// Ignore already existing import error
		collection.AddImport(name, i.engine)

		i.collectionImports[collection.Name] = NewCollectionImport(collection)
	}

	return
}

// Remove import from the list
func (c *Classify) DeleteImports(importList map[string]*Import, collections map[string]*Collection) (err error) {

	// At least one import id or one collection must be specified
	if len(importList) == 0 && len(collections) == 0 {
		err = errors.New("required import names or collection names")
		return
	}

	// Stop all imports
	c.StopImports(importList, collections)

	// If no importList are specified : remove all import relative to the
	// same collection
	if len(importList) == 0 {
		importList = c.imports
	}

	for importName, i := range importList {

		// Unlink the collection with the specified import
		for collectionName, collection := range collections {

			if err = i.Unlink2DB(c.database, collection); err != nil {
				return
			}

			// Unlink in the collection
			collection.DeleteImport(importName)

			// in the import collection list
			delete(i.collectionImports, collectionName)
		}

		// If no collection are linked with specified import
		if len(i.collectionImports) < 1 {

			if err = i.Delete2DB(c.database); err != nil {
				return
			}

			// Remove the import
			delete(c.imports, importName)
		}
	}
	return
}

// Get the whole list of imports by Type
func (c *Classify) GetImports(importList map[string]*Import, collections map[string]*Collection) (res map[string]map[string]imports.Import, err error) {

	res = make(map[string]map[string]imports.Import)

	// If no importList are specified : get all
	if len(importList) == 0 {
		importList = c.imports
	}

	for name, i := range importList {

		if i.HasCollections(collections) == false {
			continue
		}

		ref := i.engine.GetRef()

		if res[ref.String()] == nil {
			res[ref.String()] = make(map[string]imports.Import)
		}

		res[ref.String()][name] = i.engine
	}

	return
}

func (c *Classify) SendImportEvent(name string, status bool) {

	var statusStr string
	if status {
		statusStr = "start"
	} else {
		statusStr = "end"
	}

	c.SendEvent("import/status", statusStr, name, status)
}

// Launch the process of importation of specified imports
func (c *Classify) StartImports(imports map[string]*Import, collections map[string]*Collection) error {

	// If no imports are specified : get all
	if len(imports) == 0 {
		imports = c.imports
	}

	// Get the import channel
	for name, i := range imports {

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
			c.SendImportEvent(name, true)

			for {

				if input, ok := <-channel; ok {

					// For each collections linked with the importation
					for _, collectionImport := range i.collectionImports {

						collection := collectionImport.Collection

						id := data.GetId(input)
						if _, err := collection.OnInput(Id(id), input); err != nil {
							log.Printf("[%s x %d] %s\n",
								collection.Name, id, err.Error())
						}
					}

					continue
				}

				break
			}

			// Send notification to stop analysis
			c.SendImportEvent(name, false)
		}()
	}

	return nil
}

// Stop the importing process
func (c *Classify) StopImports(imports map[string]*Import, collections map[string]*Collection) error {

	// If no imports are specified : get all
	if len(imports) == 0 {
		imports = c.imports
	}

	for name, i := range imports {

		if i.HasCollections(collections) == false {
			continue
		}

		c.imports[name].engine.Stop()

		// Send notification
		go c.SendImportEvent(name, false)
	}

	return nil
}

func (c *Classify) GetImportRefs() []string {
	return imports.REF_IDX2STR
}
