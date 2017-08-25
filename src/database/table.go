package database

import (
	"fmt"
	"sort"
	"strings"
)

type Table struct {
	Name       string
	Attributes map[string]*Attribute
	Unique     []string
}

func (t *Table) HasAttribute(name string) (ok bool) {
	_, ok = t.Attributes[name]
	return
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

	// Handle unicity on multiples attributes
	if len(t.Unique) > 0 {
		result += " UNIQUE(" + strings.Join(t.Unique, ",") + ") ON CONFLICT REPLACE"

		// Remove last if ends with coma
	} else if last := len(result) - 1; last >= 0 && result[last] == ',' {

		result = result[0:last]
	}

	result += ");"
	return
}

func (t *Table) Insert(isRef bool) (result string) {

	result = "INSERT"

	if isRef {
		result += " OR IGNORE"
	}

	result += " INTO " + t.Name +
		" (" + strings.Join(t.GetAttributes(true, ""), ",") + ")" +
		" VALUES (" + strings.Join(t.GetAttributes(true, ":"), ",") + ")"

	return
}

func (t *Table) Delete(condition string) (result string) {

	result = "DELETE FROM " + t.Name + " WHERE " + condition
	return
}

func (t *Table) SelectAll() (result string) {

	result = "SELECT * FROM " + t.Name
	return
}
