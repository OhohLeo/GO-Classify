package core

import (
	"github.com/ohohleo/classify/collections"
	"github.com/ohohleo/classify/database"
	"github.com/ohohleo/classify/exports"
	"github.com/ohohleo/classify/imports"
	"log"
)

// StartDB init db and retreive all stored data
func (c *Classify) StartDB(config *Config) (err error) {

	// Activate database if enabled
	c.database, err = database.New(config.DataBase)
	if c.database == nil {
		log.Println("Database disable")
		return
	}

	log.Println("Starting Database")

	// Init database tables
	if err = collections.INIT_DB(c.database); err != nil {
		return
	}

	if err = imports.INIT_DB(c.database); err != nil {
		return
	}

	if err = exports.INIT_DB(c.database); err != nil {
		return
	}

	// Create tables
	if err = c.database.Create(); err != nil {
		return
	}

	// Insert all references
	if err = collections.INIT_REF_DB(c.database); err != nil {
		return
	}

	if err = imports.INIT_REF_DB(c.database); err != nil {
		return
	}

	if err = exports.INIT_REF_DB(c.database); err != nil {
		return
	}

	// Retreive all stored collections
	err = collections.RetreiveDBCollections(c.database,
		func(id uint64, name string, ref collections.Ref, params []byte) (err error) {
			collection, err := c.AddCollection(name, ref, params)
			if err != nil {
				return
			}

			// Store database id
			collection.Id = id
			return
		})
	if err != nil {
		return
	}

	// Retreive all stored imports
	err = imports.RetreiveDBImports(c.database,
		func(id uint64, name string, ref imports.Ref, params []byte, names []string) (err error) {

			collections, err := c.GetCollectionsByNames(names)
			if err != nil {
				return
			}

			i, _, err := c.AddImport(name, ref, params, collections)
			if err != nil {
				return
			}

			// Store database id
			i.Id = id
			return
		})
	if err != nil {
		return
	}

	// Retreive all stored exports
	err = exports.RetreiveDBExports(c.database,
		func(id uint64, name string, ref exports.Ref, params []byte, names []string) (err error) {

			collections, err := c.GetCollectionsByNames(names)
			if err != nil {
				return
			}

			e, err := c.AddExport(name, ref, params, collections)
			if err != nil {
				return
			}

			// Store database id
			e.Id = id
			return
		})
	if err != nil {
		return
	}

	return
}
