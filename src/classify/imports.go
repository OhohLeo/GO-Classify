package classify

import (
	"encoding/json"
	"errors"
	"github.com/ohohleo/classify/imports"
	"github.com/ohohleo/classify/imports/directory"
)

// Type of imports
var newImports = map[string]func(string, json.RawMessage) (imports.Import, error){
	"directory": func(name string, input json.RawMessage) (i imports.Import, err error) {
		var directory directory.Directory
		err = json.Unmarshal(input, &directory)
		if err == nil {
			directory.Name = name
			i = &directory
		}
		return
	},
}

// CreateNewImport returns the import module depending on the parameter
// given
func (c *Classify) CreateNewImport(m MappingParams) (name string, i imports.Import, err error) {

	name = getRandomName()

	if m.Type == "" {
		err = errors.New("type field is mandatory")
		return
	}

	createImport, ok := newImports[m.Type]
	if ok == false {
		err = errors.New("import type '" + m.Type + "' not handled")
		return
	}

	i, err = createImport(name, m.Params)
	if err != nil {
		return
	}

	return
}

type Import struct {
	Id          string
	engine      imports.Import
	collections []Collection
}

// func (c *Classify) AddImport() (Import, error) {

// }

// func (c *Classify) DeleteImport(id string) {

// }

// func (c *Classify) GetImports() {

// }
