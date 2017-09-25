package core

import (
	"encoding/json"
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ohohleo/classify/imports"
	"net/http"
)

// List all config imports
// GET /imports/config
func (c *Classify) ApiGetImportsConfig(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson(c.config.Imports)
}

// getImportByName get from Url parameters import
func (c *Classify) getImportByName(w rest.ResponseWriter, r *rest.Request) *Import {

	i, err := c.GetImportByName(r.PathParam("name"))
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}

	return i
}

// getImportNamesAndCollections get from Url parameters imports and the collections
func (c *Classify) getImportNamesAndCollections(r *rest.Request) (imports map[string]*Import, collections map[string]*Collection, err error) {

	// From the url query list
	values := r.URL.Query()

	var importList []string
	for _, name := range values["name"] {
		importList = append(importList, name)
	}

	// Check and get the import list
	imports, err = c.GetImportsByNames(importList)
	if err != nil {
		return
	}

	// Check and get the collection list
	collections, err = c.GetCollectionsByNames(values["collection"])
	if err != nil {
		return
	}

	return
}

type ApiAddImportsBody struct {
	Name        string          `json:"name"`
	Ref         string          `json:"ref"`
	Collections []string        `json:"collections"`
	Params      json.RawMessage `json:"params"`
}

// PostCollectionImport add a new import to the collection specified
// POST /imports
func (c *Classify) ApiAddImport(w rest.ResponseWriter, r *rest.Request) {

	// Get import parameters
	var body ApiAddImportsBody
	err := r.DecodeJsonPayload(&body)
	if err != nil {
		rest.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	// Check and get the collection list
	collections, err := c.GetCollectionsByNames(body.Collections)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Field required
	ref, ok := imports.REF_STR2IDX[body.Ref]
	if ok == false {
		rest.Error(w, fmt.Sprintf("invalid ref '%s'", body.Ref), http.StatusBadRequest)
		return
	}

	i, outParams, err := c.AddImport(body.Name, ref, body.Params, collections)
	if err != nil {

		// Manque de paramètres
		if outParams != nil {
			w.WriteJson(outParams)
			return
		}

		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Handle DB storage
	if err := i.Store2DB(c.database); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteJson(i)
}

// List all the imports selected by id or by collections
// GET /imports?name=IMPORT_NAME&collection=COLLECTION_NAME
func (c *Classify) ApiGetImports(w rest.ResponseWriter, r *rest.Request) {

	names, collections, err := c.getImportNamesAndCollections(r)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := c.GetImports(names, collections)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteJson(res)
}

// ApiDeleteImport remove specified import selected by id and by the
// collections
// DELETE /imports?name=IMPORT_NAME&collection=COLLECTION_NAME
func (c *Classify) ApiDeleteImport(w rest.ResponseWriter, r *rest.Request) {

	names, collections, err := c.getImportNamesAndCollections(r)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := c.DeleteImports(names, collections); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Launch the analysis of the collection import
// PUT /imports/start?name=IMPORT_NAME&collection=COLLECTION_NAME
func (c *Classify) ApiStartImport(w rest.ResponseWriter, r *rest.Request) {

	names, collections, err := c.getImportNamesAndCollections(r)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := c.StartImports(names, collections); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Stop the analysis of the collection import
// PUT /imports/stop?name=IMPORT_NAME&collection=COLLECTION_NAME
func (c *Classify) ApiStopImport(w rest.ResponseWriter, r *rest.Request) {

	names, collections, err := c.getImportNamesAndCollections(r)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := c.StopImports(names, collections); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Handle import params
// PUT /imports/:name/:param
func (c *Classify) ApiPutImportParam(w rest.ResponseWriter, r *rest.Request) {

	param := r.PathParam("param")
	name := r.PathParam("name")

	// Récupération du type de l'importation
	i, err := c.GetImportByName(name)
	if err == nil {
		name = i.engine.GetRef().String()
	}

	newImport, ok := newImports[name]
	if ok == false {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var body json.RawMessage
	err = r.DecodeJsonPayload(&body)
	if err != nil {
		rest.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	res, err := newImport.GetParam(param, body)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteJson(res)
	return
}
