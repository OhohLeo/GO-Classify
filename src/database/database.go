package database

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"sort"
	"strings"
)

type Config struct {
	Enable bool   `json:"enable"`
	Driver string `json:"driver"`
	Source string `json:"source"`
}

func (c *Config) Connect() (*sqlx.DB, error) {
	db, err := sqlx.Connect(c.Driver, c.Source)
	if err != nil {
		return nil, err
	}

	return db, nil
}

type RequireDB interface {
	GetDBTables() []*Table
}

type HandleDB interface {
	GetDBAttributes() map[string]interface{}
}

func Create(db *sqlx.DB, rDB RequireDB) error {

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	for _, table := range rDB.GetDBTables() {
		query, err := table.Create()
		if err != nil {
			return err
		}

		log.Info("DB:" + query)

		_, err = tx.Exec(query)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func Insert(db *sqlx.DB, t *Table, add HandleDB) error {

	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	query, err := t.Insert(false)
	if err != nil {
		return err
	}

	log.Info("DB:" + query)

	_, err = tx.NamedExec(query, add.GetDBAttributes())
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func InsertRef(db *sqlx.DB, t *Table, refs []string) error {

	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	query, err := t.Insert(true)
	if err != nil {
		return err
	}

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

type AttributeType int

const (
	TEXT AttributeType = iota
	INTEGER
)

var attributeType2str = []string{
	"TEXT",
	"INTEGER",
}

type Attribute struct {
	Type         AttributeType
	IsNotNull    bool
	IsPrimaryKey bool
	IsUnique     bool
}

func (a *Attribute) Create() string {
	res := attributeType2str[a.Type]

	if a.IsNotNull {
		res += " NOT NULL"
	}

	if a.IsPrimaryKey {
		res += " PRIMARY KEY"
	}

	if a.IsUnique {
		res += " UNIQUE"
	}

	return res
}

func (a *Attribute) String() string {
	return attributeType2str[a.Type]
}

type Table struct {
	Name       string
	Attributes map[string]*Attribute
}

func (t *Table) GetAttributes(ignore bool, prefix string) []string {

	var attributes []string
	idx := 0

	hasPrefix := false
	if prefix != "" {
		hasPrefix = true
	}

	for name, _ := range t.Attributes {

		if ignore && name == "id" {
			continue
		}

		if hasPrefix {
			name = prefix + name
		}

		attributes = append(attributes, name)
		idx++
	}

	// Alphabetic sort
	sort.Strings(attributes)

	return attributes
}

func (t *Table) Create() (result string, err error) {

	if t.Name == "" {
		err = fmt.Errorf("table should not have empty name")
		return
	}

	result = "CREATE TABLE IF NOT EXISTS " + t.Name + " ("

	for _, name := range t.GetAttributes(false, "") {
		result += name + " " + t.Attributes[name].Create() + ","
	}

	// Remove last if ends with coma
	if last := len(result) - 1; last >= 0 && result[last] == ',' {
		result = result[0:last]
	}

	result += ");"
	return
}

func (t *Table) Insert(isRef bool) (result string, err error) {

	result = "INSERT"

	if isRef {
		result += " OR IGNORE"
	}

	result += " INTO " + t.Name +
		" (" + strings.Join(t.GetAttributes(true, ""), ",") + ")" +
		" VALUES (" + strings.Join(t.GetAttributes(true, ":"), ",") + ")"

	return
}
