package collections

import (
	"github.com/jmoiron/sqlx"
	"github.com/ohohleo/classify/database"
)

var DB_DETAILS database.Table = database.Table{
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

func (c *Collection) GetDBAttributes() map[string]interface{} {
	return map[string]interface{}{
		"name":   c.name,
		"type":   c.GetType(),
		"params": "",
	}
}

var DB_TYPES_REF database.Table = database.Table{
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
	Params string `db:"params"`
}

func GetDBCollections(db *sqlx.DB) (collections []DB, err error) {

	err = db.Select(&collections, "SELECT collections.* FROM collections")
	return
}
