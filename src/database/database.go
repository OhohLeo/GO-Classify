package database

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"regexp"
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

var reTableName = regexp.MustCompile("([a-z0-9_]+)_id$")

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

	var unique []string

	// Initialize attributes
	attributes := make(map[string]*Attribute)
	for _, name := range parameters {

		var ok bool
		var attribute *Attribute

		// Handle specific {table_name}_id attributes
		if tableNameId := reTableName.FindStringSubmatch(name); len(tableNameId) > 1 {

			tableId, ok := d.tables[tableNameId[1]]
			if ok == false {
				return fmt.Errorf("Table '%s' not found and referenced by '%s'",
					tableNameId, table)
			}

			// Check referenced table exists
			if tableId.HasAttribute("id") == false {
				return fmt.Errorf("No attribute id found on table '%s' referenced by '%s'",
					tableNameId, table)
			}

			unique = append(unique, name)
			attribute = &Attribute{Type: INTEGER}

		} else {

			// Handle generic attributes
			attribute, ok = genAttributes[name]
			if ok == false {
				return fmt.Errorf("no DB generic attribute found for table '%s/%s'",
					table, name)
			}
		}

		attributes[name] = attribute
	}

	// Store table
	d.tables[table] = &Table{
		Name:       table,
		Attributes: attributes,
		Unique:     unique,
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

		log.Println("DB [" + name + "] " + query)

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

	log.Println("DB [" + name + "] " + query + fmt.Sprintf(" [%+v]", toStore))

	_, err = tx.NamedExec(query, toStore)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (d *Database) Insert(name string, toStore interface{}) (uint64, error) {

	table, err := d.GetTable(name)
	if err != nil {
		return 0, err
	}

	tx, err := d.db.Beginx()
	if err != nil {
		return 0, err
	}

	query := table.Insert(false)

	log.Println("DB [" + name + "] " + query)

	row, err := tx.NamedExec(query, toStore)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	id, err := row.LastInsertId()
	return uint64(id), err
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

	log.Println("DB [" + name + "] " + query)

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

func (d *Database) Select(result interface{}, query string, params ...interface{}) (err error) {

	log.Println("DB " + query)

	err = d.db.Select(result, query, params...)
	return
}

func (d *Database) SelectAll(name string) (result []GenStruct, err error) {

	table, err := d.GetTable(name)
	if err != nil {
		return
	}

	query := table.SelectAll()

	log.Println("DB [" + name + "] " + query)

	err = d.db.Select(&result, query)
	return
}
