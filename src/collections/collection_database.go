package collections

import (
	"encoding/json"
	"github.com/ohohleo/classify/database"
)

func INIT_DB(db *database.Database) (err error) {

	err = db.AddTable("collections", []string{
		"id", "name", "ref", "params"})
	if err != nil {
		return
	}

	err = db.AddTable("collections_refs", []string{
		"id", "name"})
	if err != nil {
		return
	}

	return
}

func INIT_REF_DB(db *database.Database) error {
	return db.InsertRef("collections_refs", REF_IDX2STR)
}

type Params struct {
	Websites []string `json:"websites"`
}

func (c *Collection) Store2DB(db *database.Database) error {

	// Store websites data
	params := &Params{
		Websites: c.GetWebsiteParams(),
	}

	// Convert to JSON
	paramsStr, err := json.Marshal(params)
	if err != nil {
		return err
	}

	// Store the collection
	return db.Insert("collections", &database.GenStruct{
		Name:   c.name,
		Ref:    uint64(c.GetRef()),
		Params: paramsStr,
	})
}

func (c *Collection) Delete2DB(db *database.Database) error {

	return db.Delete("collections", &database.GenStruct{
		Name: c.name,
		Ref:  uint64(c.GetRef()),
	}, "name = :name AND ref = :ref")
}

type OnCollection func(name string, ref Ref, params Params) error

func RetreiveDBCollections(db *database.Database, onCollection OnCollection) (err error) {

	var dbCollections []database.GenStruct
	dbCollections, err = db.SelectAll("collections")
	if err != nil {
		return
	}

	for _, dbCollection := range dbCollections {

		// Get collection params
		var params Params
		err = dbCollection.GetParams(&params)
		if err != nil {
			return
		}

		// Add new stored collection
		err = onCollection(dbCollection.Name, Ref(dbCollection.Ref), params)
		if err != nil {
			return
		}
	}

	return
}
