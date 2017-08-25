package database

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type GenStruct struct {
	Id     uint64 `db:"id"`
	Name   string `db:"name"`
	Ref    uint64 `db:"ref"`
	Params []byte `db:"params"`
}

func (stored *GenStruct) GetParams(params interface{}) error {
	return json.Unmarshal(stored.Params, params)
}

var genAttributes = map[string]*Attribute{
	"id": &Attribute{
		Type:         INTEGER,
		IsPrimaryKey: true,
	},
	"name": &Attribute{
		Type:     TEXT,
		IsUnique: true,
	},
	"ref": &Attribute{
		Type: INTEGER,
	},
	"params": &Attribute{
		Type: TEXT,
	},
}

type Database struct {
	db     *sqlx.DB
	tables map[string]*Table
}

type Config struct {
	Enable bool   `json:"enable"`
	Driver string `json:"driver"`
	Source string `json:"source"`
}

func New(config Config) (d *Database, err error) {

	// Check if database is enable
	if config.Enable == false {
		return
	}

	d = new(Database)

	// Establish database connection
	d.db, err = sqlx.Connect(config.Driver, config.Source)
	if err != nil {
		return
	}

	return
}

func (d *Database) AddTable(table string, parameters []string) error {

	// Check if table name already exists
	if _, ok := d.tables[table]; ok {
		return fmt.Errorf("already existing DB table '%s'", table)
	}

	if len(parameters) == 0 {
		return fmt.Errorf("no DB parameters found for table '%s'", table)
	}

	// Initialize tables if necessary
	if d.tables == nil {
		d.tables = make(map[string]*Table)
	}

	// Initialize attributes
	attributes := make(map[string]*Attribute)
	for _, name := range parameters {

		attribute, ok := genAttributes[name]
		if ok == false {
			return fmt.Errorf("no DB generic attribute found for table '%s/%s'",
				table, name)
		}

		attributes[name] = attribute
	}

	// Store table
	d.tables[table] = &Table{
		Name:       table,
		Attributes: attributes,
	}

	return nil
}

func (d *Database) GetTable(name string) (table *Table, err error) {

	var ok bool
	table, ok = d.tables[name]
	if ok == false {
		err = fmt.Errorf("no DB table '%s' found", table)
		return
	}

	return
}

func (d *Database) Create() error {

	tx, err := d.db.Begin()
	if err != nil {
		return err
	}

	for name, table := range d.tables {
		query, err := table.Create()
		if err != nil {
			return err
		}

		log.Info("DB " + name + " : " + query)

		_, err = tx.Exec(query)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (d *Database) Delete(name string, toStore interface{}, condition string) error {

	table, err := d.GetTable(name)
	if err != nil {
		return err
	}

	tx, err := d.db.Beginx()
	if err != nil {
		return err
	}

	query := table.Delete(condition)

	log.Info("DB " + name + " " + query)

	_, err = tx.NamedExec(query, toStore)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// func (SimpleDB *s) Store(db *sqlx.DB, t *Table, params interface{}) error {

// 	// Convert to JSON
// 	s.Params, err := json.Marshal(params)
// 	if err != nil {
// 		return err
// 	}

// 	// Store the collection
// 	return Insert(db, &DB_LIST, s)
// }

func (d *Database) Insert(name string, toStore interface{}) error {

	table, err := d.GetTable(name)
	if err != nil {
		return err
	}

	tx, err := d.db.Beginx()
	if err != nil {
		return err
	}

	query := table.Insert(false)

	log.Info("DB:" + query)

	_, err = tx.NamedExec(query, toStore)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (d *Database) InsertRef(name string, refs []string) error {

	table, err := d.GetTable(name)
	if err != nil {
		return err
	}

	tx, err := d.db.Beginx()
	if err != nil {
		return err
	}

	query := table.Insert(true)

	log.Info("DB:" + query)

	for _, ref := range refs {

		if _, err = tx.NamedExec(query, map[string]interface{}{
			"name": ref,
		}); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (d *Database) SelectAll(name string) (result []GenStruct, err error) {

	table, err := d.GetTable(name)
	if err != nil {
		return
	}

	query := table.SelectAll()

	log.Info("DB:" + query)

	err = d.db.Select(&result, query)
	return
}
