package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ohohleo/classify/collections"
	"github.com/ohohleo/classify/data"
	"github.com/ohohleo/classify/database"
	"github.com/ohohleo/classify/exports"
	"github.com/ohohleo/classify/imports"
	"github.com/ohohleo/classify/websites"
	"log"
)

type Collection struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`

	buffer *Buffer
	items  *Items
	config *CollectionConfig

	events chan CollectionEvent
	engine collections.Collection

	imports  map[string]imports.Import
	exports  []exports.Export
	websites map[string]websites.Website
}

type CollectionEvent struct {
	Source string
	Status string
	Id     string
	Item   *Item
}

type CollectionParams struct {
	Websites []string
}

// Add new import to the collection
func (c *Collection) AddImport(name string, i imports.Import) error {

	if c.imports == nil {
		c.imports = make(map[string]imports.Import)
	}

	if _, ok := c.imports[name]; ok {

		return fmt.Errorf("collection '%s' has already import '%s'",
			c.Name, name)
	}

	c.imports[name] = i

	// Update data configuration
	c.config.UpdateDatas(i)

	return nil
}

// Remove existing import from the collection
func (c *Collection) DeleteImport(name string) error {

	if _, ok := c.imports[name]; ok == false {

		return fmt.Errorf("import '%s' not found on collection '%s'",
			name, c.Name)
	}

	delete(c.imports, name)

	// Reset datas config
	c.config.Datas = nil

	// Refresh new list
	for _, i := range c.imports {
		c.config.UpdateDatas(i)
	}

	return nil
}

func (c *Collection) ActivateStore() {
}

func (c *Collection) DisableStore() error {
	return nil
}

func (c *Collection) Store2DB(db *database.Database) error {

	// Store websites data
	params := &CollectionParams{
		Websites: c.GetWebsiteParams(),
	}

	// Convert to JSON
	configJson, err := json.Marshal(c.config)
	if err != nil {
		return err
	}

	// Convert to JSON
	paramsJson, err := json.Marshal(params)
	if err != nil {
		return err
	}

	// Store the collection
	lastId, err := db.Insert("collections", &database.GenStruct{
		Name:   c.Name,
		Ref:    uint64(c.engine.GetRef()),
		Config: configJson,
		Params: paramsJson,
	})

	c.Id = lastId

	return err
}

// Store configuration on DataBase
func (c *Collection) StoreConfig2DB(db *database.Database) error {
	return c.config.Store2DB(c, db)
}

func (c *Collection) Delete2DB(db *database.Database) error {

	return db.Delete("collections", &database.GenStruct{
		Name: c.Name,
		Ref:  uint64(c.engine.GetRef()),
	}, "name = :name AND ref = :ref")
}

// AddWebsite add new website
func (c *Collection) AddWebsite(website websites.Website) {

	if c.websites == nil {
		c.websites = make(map[string]websites.Website)
	}

	c.websites[website.GetRef().String()] = website
}

// DeleteWebsite delete specified website
func (c *Collection) DeleteWebsite(name string) error {

	if _, ok := c.websites[name]; ok {

		delete(c.websites, name)
		return nil
	}

	return errors.New("no website name '" + name + "' found")
}

func (c *Collection) GetWebsiteParams() []string {

	keys := make([]string, len(c.websites))

	idx := 0
	for name := range c.websites {
		keys[idx] = name
		idx++
	}

	return keys
}

func (c *Collection) ActivateBuffer() {

	// Init the items buffer
	c.ResetBuffer()
}

func (c *Collection) DisableBuffer() error {

	if c.buffer == nil {
		return fmt.Errorf("collection '%s' buffer already disabled", c.Name)
	}

	c.buffer.RemoveAll()
	c.buffer = nil
	return nil
}

func (c *Collection) Search(src string, item *BufferItem) {

	// Launch research through web
	if len(c.websites) > 0 {
		c.SearchWeb(item)

		// item.MatchId = item.Websites["TMDB"][0].GetId()
	}

	// // TODO Research for best matching
	// item.Ref = item.Websites["IMDB"][0].GetRef()
	// item.IsMatching = 10.2
	// item.BestMatch = item.Websites["IMDB"][0]

	c.SendCollectionEvent(src, "update", &item.Item)
}

// SearchWeb launch resarch through specified websites
func (c *Collection) SearchWeb(item *BufferItem) {

	// Get name to search
	keywords := item.WebQuery

	// For all specified websites
	for _, website := range c.websites {

		// Launch the research
		channel := website.Search(keywords)

		for {
			d, ok := <-channel
			if ok {
				item.AddWebsiteData(website.GetRef().String(), d)
				continue
			}

			break
		}
	}
}

// OnInput handle new data to classify
func (c *Collection) OnInput(input data.Data) (*BufferItem, error) {

	if err := c.onDataInput(input); err != nil {
		log.Printf("OnDataInput error: %+v\n", err)
	}

	// Create a new item
	item := NewBufferItem()

	// log.Printf("OnInput %+v\n", input)

	if c.buffer != nil {

		// Add the import to the item
		item.AddImportData(input)

		// Get Cleaned name
		item.SetCleanedName(c.config.Import.Banned, c.config.Import.Separators)

		// Store item to the buffer collection
		c.buffer.Add(item.Item.Engine.GetName(), item)

	} else {

		item.Item.Id = getRandomId()
		item.Item.SetData(input)

		// Otherwise directly store item to the items collection
		err := c.items.Add(item.Item.Engine.GetName(), &item.Item)
		if err != nil {
			return nil, err
		}

		c.SendCollectionEvent("items", "add", &item.Item)
	}

	return item, nil
}

func (c *Collection) onDataInput(d data.Data) error {

	// Get data ref name
	refName := d.GetRef().String()

	// Call method OnCollection if it exists
	if inputData, ok := d.(data.OnCollection); ok {
		if err := inputData.OnCollection(c.config.Datas[refName]); err != nil {
			return err
		}
	}

	// Check for data dependencies
	if hasDeps, ok := d.(data.HasDependencies); ok {
		for _, dep := range hasDeps.GetDependencies() {
			if err := c.onDataInput(dep); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Collection) SendCollectionEvent(src string, status string, item *Item) {

	if c.events == nil {
		return
	}

	c.events <- CollectionEvent{
		Source: src,
		Status: status,
		Id:     item.Engine.GetName(),
		Item:   item,
	}
}

func (c *Collection) GetBuffer() []*BufferItem {

	if c.buffer == nil {
		return []*BufferItem{}
	}

	return c.buffer.GetCurrentList()
}

func (c *Collection) ResetBuffer() {
	c.buffer = NewBuffer(c, 2)
}

func (c *Collection) GetBufferItems() []*Item {

	if c.items == nil {
		return []*Item{}
	}

	return c.items.GetCurrentList()
}

func (c *Collection) ResetItems() {
	c.items = NewItems()
}

func (c *Collection) Validate(id string, d data.Data) (item *BufferItem, err error) {

	if c.buffer == nil {
		err = fmt.Errorf("buffer not initialized")
		return
	}

	item = c.buffer.Validate(id)
	if item == nil {
		err = fmt.Errorf("name '%s' not referenced in buffer", id)
		return
	}

	// Store in definitive items
	err = c.items.Add(id, &item.Item)
	if err != nil {
		return
	}

	c.SendCollectionEvent("items", "add", &item.Item)

	// TODO: export list

	return
}

func (c *Collection) GetItems() []*Item {

	return c.items.GetCurrentList()
}

func (c *Collection) SetExports(exports []exports.Export) {
	c.exports = exports
}

// OnOutput export data from classify
func (c *Collection) onOutput(item *BufferItem) error {

	return nil
}
