package collections

import (
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"github.com/ohohleo/classify/database"
)

var DB_LIST database.Table = database.Table{
	Name: "collections",
	Attributes: map[string]*database.Attribute{
		"id": &database.Attribute{
			Type:         database.INTEGER,
			IsPrimaryKey: true,
		},
		"name": &database.Attribute{
			Type:     database.TEXT,
			IsUnique: true,
		},
		"type": &database.Attribute{
			Type: database.INTEGER,
		},
		"params": &database.Attribute{
			Type: database.TEXT,
		},
	},
}

var DB_REFS database.Table = database.Table{
	Name: "type_refs",
	Attributes: map[string]*database.Attribute{
		"id": &database.Attribute{
			Type:         database.INTEGER,
			IsPrimaryKey: true,
		},
		"name": &database.Attribute{
			Type:     database.TEXT,
			IsUnique: true,
		},
	},
}

type DB struct {
	Id     uint64 `db:"id"`
	Name   string `db:"name"`
	Type   uint64 `db:"type"`
	Params []byte `db:"params"`
}

func (stored *DB) GetParams() (params Params, err error) {
	err = json.Unmarshal(stored.Params, &params)
	return
}

type Params struct {
	Websites []string `json:"websites"`
}

func (c *Collection) Store2DB(db *sqlx.DB) error {

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
	return database.Insert(db, &DB_LIST, &DB{
		Name:   c.name,
		Type:   uint64(c.GetType()),
		Params: paramsStr,
	})
}

func GetDBCollections(db *sqlx.DB) (collections []DB, err error) {

	err = db.Select(&collections, "SELECT collections.* FROM collections")
	return
}
