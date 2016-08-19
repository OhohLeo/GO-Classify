package classify

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ohohleo/classify/imports"
	"net/http"
)

// PostCollectionImport add a new import to the collection specified
// POST /collection/:name/imports
func (c *Classify) PostCollectionImport(w rest.ResponseWriter, r *rest.Request) {

	// Check the collection exist
	collection := c.getCollectionByName(w, r)
	if collection == nil {
		return
	}

	// Get mapping
	var m MappingParams
	err := r.DecodeJsonPayload(&m)
	if err != nil {
		rest.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	// Create the new import
	name, i, err := c.CreateNewImport(m)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = collection.AddImport(name, i)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Add("References", name)

	w.WriteHeader(http.StatusNoContent)
}

// GetCollectionImportByName returns the existing collection import specified
func (c *Classify) GetCollectionImportByName(w rest.ResponseWriter, r *rest.Request) imports.Import {

	collection, err := c.GetCollection(r.PathParam("name"))
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}

	imported, err := collection.GetImportByName(r.PathParam("import"))
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}

	return imported
}

// GetCollectionImport list all the imports used by the collection specified
// GET /collection/:name/imports
func (c *Classify) GetCollectionImports(w rest.ResponseWriter, r *rest.Request) {

	// Check the collection exist
	collection := c.getCollectionByName(w, r)
	if collection == nil {
		return
	}

	w.WriteJson(collection.GetImports())
}

// DeleteCollectionImport remove specified imports used by the collection specified
// DELETE /collection/:name/imports/:import
func (c *Classify) DeleteCollectionImport(w rest.ResponseWriter, r *rest.Request) {

	// Check the collection exist
	collection := c.getCollectionByName(w, r)
	if collection == nil {
		return
	}

	if err := collection.DeleteImport(r.PathParam("import")); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// StartCollectionImport launch the analysis of the collection import
// PUT /collection/:name/imports/:import/start
func (c *Classify) StartCollectionImport(w rest.ResponseWriter, r *rest.Request) {

	// Check the imported exist
	imported := c.GetCollectionImportByName(w, r)
	if imported == nil {
		return
	}

	channel, err := imported.Start()
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	go func() {
		for {
			if item, ok := <-channel; ok {
				c.Server.Send(r.PathParam("name"), "newFile", item)
				continue
			}
			break
		}
	}()

	w.WriteHeader(http.StatusNoContent)
}

// StopCollectionImport stop the analysis of the collection import
// PUT /collection/:name/imports/:import/stop
func (c *Classify) StopCollectionImport(w rest.ResponseWriter, r *rest.Request) {

	// Check the imported exist
	imported := c.GetCollectionImportByName(w, r)
	if imported == nil {
		return
	}

	imported.Stop()

	w.WriteHeader(http.StatusNoContent)
}
