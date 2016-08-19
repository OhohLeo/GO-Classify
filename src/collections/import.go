package collections

import (
	"errors"
	"github.com/ohohleo/classify/imports"
)

// AddImport add new import
func (c *Collection) AddImport(name string, imported imports.Import) error {

	if _, ok := c.imports[name]; ok {
		return errors.New("import '" + name + "' already exists")
	}

	if c.imports == nil {
		c.imports = make(map[string]imports.Import)
	}

	c.imports[name] = imported

	if c.items != nil {
		c.startImport(imported)
	}

	return nil
}

// DeleteImport delete specified import
func (c *Collection) DeleteImport(name string) error {

	imported, err := c.GetImportByName(name)
	if err != nil {
		return err
	}

	imported.Stop()
	delete(c.imports, name)
	return nil
}

// GetImports get the list of imports
func (c *Collection) GetImports() map[string]map[string]imports.Import {

	result := make(map[string]map[string]imports.Import)

	if c.imports == nil {
		return result
	}

	for name, imported := range c.imports {

		t := imported.GetType()

		if result[t] == nil {
			result[t] = make(map[string]imports.Import)
		}

		result[t][name] = imported
	}

	return result
}

// Returns an error if import doesn't exist on the specified collection
func (c *Collection) GetImportByName(name string) (imports.Import, error) {

	// If imports are initialised
	if c.imports != nil {

		// And the import name exists
		imported, ok := c.imports[name]
		if ok {

			// It is OK
			return imported, nil
		}
	}

	return nil, errors.New("import '" + name + "' not found")
}

// Start the analysis of the import specified
func (c *Collection) StartSpecific(name string) (chan *Item, error) {

	imported, err := c.GetImportByName(name)
	if err != nil {
		return nil, err
	}

	items := c.initItems()

	c.startImport(imported)

	return items, nil
}

// Stop the analysis of the import specified
func (c *Collection) StopSpecific(name string) error {

	imported, err := c.GetImportByName(name)
	if err != nil {
		return err
	}

	imported.Stop()
	return nil
}

// startImports launch the process of importation of specified import
func (c *Collection) startImport(imported imports.Import) error {

	// Get the import channel
	channel, err := imported.Start()
	if err != nil {
		return err
	}

	// Send all data imported to the collection
	go func() {
		for {
			if input, ok := <-channel; ok {
				c.items <- c.OnInput(input)
				continue
			}
			break
		}
	}()

	return nil
}
