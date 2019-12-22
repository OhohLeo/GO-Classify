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
	"github.com/ohohleo/classify/params"
	"github.com/ohohleo/classify/reference"
)

// Type of imports
var newImports = map[string]imports.Build{
	"directory": directory.ToBuild(),
	"imap":      imap.ToBuild(),
}

func Import2Build(typ string) (imports.Build, error) {
	buildImport, ok := newImports[typ]
	if ok == false {
		return buildImport, fmt.Errorf("import type '%s' not handled", typ)
	}
	return buildImport, nil
}

type Import struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`

	engine imports.Import

	configs map[*Collection]*Configs
	params  map[string]params.Param
}

func NewImport(typ string) (*Import, error) {
	buildImport, err := Import2Build(typ)
	if err != nil {
		return nil, err
	}

	return &Import{
		Name:   typ,
		engine: buildImport.ForceCreate(),
	}, nil
}

func (i *Import) GetEngine() imports.Import {
	return i.engine
}

func (i *Import) GetRefs() []*reference.Ref {
	return reference.GetRefs(i.engine)
}

func (i *Import) GetDatas() map[string]interface{} {
	result := make(map[string]interface{})
	for _, data := range i.engine.GetDatasReferences() {
		result[data.GetRef().String()] = data
	}
	return result
}

func (i *Import) GetDatasReferences() DatasReference {
	return GetDatasReference(i.GetDatas())
}

func (i *Import) HasCollection(collection *Collection) (ok bool) {
	_, ok = i.configs[collection]
	return
}

// Return true if import has a specified collection or no collections are specified
func (i *Import) HasCollections(collections map[string]*Collection) bool {

	if len(collections) > 0 {

		for _, collection := range collections {
			if i.HasCollection(collection) {
				return true
			}
		}

		// No collection match
		return false
	}

	return true
}

func (i *Import) GetConfig(collection *Collection) (configs *Configs, err error) {

	var ok bool
	if configs, ok = i.configs[collection]; !ok {
		err = fmt.Errorf("no config found for collection '%s'", collection.Name)
	}
	return
}

func (i *Import) SetConfig(collection *Collection, newConfigs *Configs) (err error) {

	var configs *Configs
	configs, err = i.GetConfig(collection)
	if err != nil {
		return
	}

	if newConfigs.Generic != nil {

		// TODO : handle enable/disable import
		if configs.Generic.Enabled != newConfigs.Generic.Enabled {
			configs.Generic.Enabled = newConfigs.Generic.Enabled
		}
	}

	if newConfigs.Tweak != nil {
		fmt.Printf("HAS TWEAK %+v\n", newConfigs.Tweak)
		configs.Tweak = newConfigs.Tweak
	}

	return
}

func (i *Import) GetParam(name string, input json.RawMessage) (interface{}, error) {

	if i.params == nil {
		i.params = params.GetParams(i.engine)
	}

	param, ok := i.params[name]
	if ok == false {
		return nil, fmt.Errorf("import param '%s' not found for '%s'", name, i.Name)
	}

	return param.ExecuteParam(input)
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
	for _, config := range i.configs {

		_, err := db.Insert("imports_mappings",
			map[string]interface{}{
				"imports_id":     lastId,
				"collections_id": config.Collection.Id,
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
func (c *Classify) CheckImportsConfig(configuration map[string]json.RawMessage) error {

	// For all import configuration
	for importType, config := range configuration {

		// Select only generic type (ie ':type_name')
		if len(importType) < 1 || importType[0] != ':' {
			continue
		}

		// Check that the import type does exists
		buildImport, err := Import2Build(importType)
		if err != nil {
			return err
		}

		// Check specified configuration
		err = buildImport.CheckConfig(config)
		if err != nil {
			return err
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

// Create import process and store it
func (c *Classify) CreateImport(name string, ref imports.Ref, inParams json.RawMessage, collections map[string]*Collection) (i *Import, outParams interface{}, err error) {
	i, outParams, err = c.AddImport(name, ref, inParams, collections)
	if err != nil {
		return
	}

	// Handle DB storage
	if err = i.Store2DB(c.database); err != nil {
		return
	}

	return
}

// Add import process
func (c *Classify) AddImport(name string, ref imports.Ref, inParams json.RawMessage, collections map[string]*Collection) (i *Import, outParams interface{}, err error) {

	// Nécessite l'existence d'au moins une collection
	if len(collections) < 1 {
		err = errors.New("required at least one existing collection")
		return
	}

	// Check that the type exists
	var buildImport imports.Build
	buildImport, err = Import2Build(ref.String())
	if err != nil {
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
	i.configs = make(map[*Collection]*Configs)

	// Add import to the collection
	for _, collection := range collections {

		// Ignore already existing import error
		collection.AddImport(name, i)

		i.configs[collection] = NewConfigs(collection, nil)
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
		for _, collection := range collections {

			if err = i.Unlink2DB(c.database, collection); err != nil {
				return
			}

			// Unlink in the collection
			collection.DeleteImport(importName)

			// in the import collection list
			delete(i.configs, collection)
		}

		// If no collection are linked with specified import
		if len(i.configs) < 1 {

			if err = i.Delete2DB(c.database); err != nil {
				return
			}

			// Remove the import
			delete(c.imports, importName)
		}
	}
	return
}

// Get the whole list of imports
func (c *Classify) GetImports(importList map[string]*Import, collections map[string]*Collection) (res map[string]imports.Import, err error) {
	res = make(map[string]imports.Import)

	// If no importList are specified : get all
	if len(importList) == 0 {
		importList = c.imports
	}

	for name, i := range importList {
		if i.HasCollections(collections) == false {
			continue
		}

		res[name] = i.engine
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
					for _, config := range i.configs {
						collection := config.Collection
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

func (c *Classify) GetImportRefs() map[string]References {
	references := make(map[string]References)

	for _, name := range imports.REF_IDX2STR {
		i, _ := NewImport(name)
		references[name] = References{
			Datas: i.GetDatasReferences(),
		}
	}

	return references
}
