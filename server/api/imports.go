package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ohohleo/classify/core"
	"github.com/ohohleo/classify/data"
	"github.com/ohohleo/classify/imports"
)

// getImportByName get from Url parameters import
func (a *API) getImportByName(w rest.ResponseWriter, r *rest.Request) *core.Import {

	i, err := a.Classify.GetImportByName(r.PathParam("name"))
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}

	return i
}

// getImportNamesAndCollections get from Url parameters imports and the collections
func (a *API) getImportNamesAndCollections(r *rest.Request) (imports map[string]*core.Import, collections map[string]*core.Collection, err error) {

	// From the url query list
	values := r.URL.Query()

	var importList []string
	for _, name := range values["name"] {
		importList = append(importList, name)
	}

	// Check and get the import list
	imports, err = a.Classify.GetImportsByNames(importList)
	if err != nil {
		return
	}

	// Check and get the collection list
	collections, err = a.Classify.GetCollectionsByNames(values["collection"])
	if err != nil {
		return
	}

	return
}

// getImportNamesAndCollections get from Url parameters imports and the collections
func (a *API) getSingleCollectionByQuery(w rest.ResponseWriter, r *rest.Request) *core.Collection {

	// From the url query list
	values := r.URL.Query()

	if len(values["collection"]) != 1 {
		rest.Error(w, "one (and only one) collection expected", http.StatusBadRequest)
		return nil
	}

	// Check and get the collection list
	collection, err := a.Classify.GetCollection(values["collection"][0])
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}

	return collection
}

type AddImportsBody struct {
	Name        string          `json:"name"`
	Ref         string          `json:"ref"`
	Collections []string        `json:"collections"`
	Params      json.RawMessage `json:"params"`
}

// PostCollectionImport add a new import to the collection specified
// POST /imports
func (a *API) AddImport(w rest.ResponseWriter, r *rest.Request) {

	// Get import parameters
	var body AddImportsBody
	err := r.DecodeJsonPayload(&body)
	if err != nil {
		rest.Error(w, fmt.Sprintf("invalid json body: %s", err.Error()), http.StatusBadRequest)
		return
	}

	// Check and get the collection list
	collections, err := a.Classify.GetCollectionsByNames(body.Collections)
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

	i, outParams, err := a.Classify.CreateImport(body.Name, ref, body.Params, collections)
	if err != nil {

		// Manque de paramètres
		if outParams != nil {
			w.WriteJson(outParams)
			return
		}

		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteJson(i)
}

// List all the imports selected by id or by collections
// GET /imports?name=IMPORT_NAME&collection=COLLECTION_NAME
func (a *API) GetImports(w rest.ResponseWriter, r *rest.Request) {
	names, collections, err := a.getImportNamesAndCollections(r)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	importList, err := a.Classify.GetImports(names, collections)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	type ImportBody struct {
		Params imports.Import `json:"params"`
		Ref    string         `json:"ref"`
	}

	res := make(map[string]ImportBody)
	for name, i := range importList {
		res[name] = ImportBody{i, i.GetRef().String()}
	}

	w.WriteJson(res)
}

// DeleteImport remove specified import selected by id and by the
// collections
// DELETE /imports?name=IMPORT_NAME&collection=COLLECTION_NAME
func (a *API) DeleteImport(w rest.ResponseWriter, r *rest.Request) {
	names, collections, err := a.getImportNamesAndCollections(r)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := a.Classify.DeleteImports(names, collections); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Launch the analysis of the collection import
// PUT /imports/start?name=IMPORT_NAME&collection=COLLECTION_NAME
func (a *API) StartImport(w rest.ResponseWriter, r *rest.Request) {
	names, collections, err := a.getImportNamesAndCollections(r)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := a.Classify.StartImports(names, collections); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Stop the analysis of the collection import
// PUT /imports/stop?name=IMPORT_NAME&collection=COLLECTION_NAME
func (a *API) StopImport(w rest.ResponseWriter, r *rest.Request) {
	names, collections, err := a.getImportNamesAndCollections(r)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := a.Classify.StopImports(names, collections); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Get import informations
// GET /imports/:name?references
func (a *API) GetImport(w rest.ResponseWriter, r *rest.Request) {
	i := a.getImportByName(w, r)
	if i == nil {
		return
	}

	collection := a.getSingleCollectionByQuery(w, r)
	if collection == nil {
		return
	}

	config, err := i.GetConfig(collection)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	type ImportBody struct {
		Ref          string         `json:"ref"`
		Input        imports.Import `json:"input"`
		OutputFormat []data.Data    `json:"output_format,omitempty"`
		Config       *core.Configs  `json:"config"`
	}

	engine := i.GetEngine()
	body := ImportBody{
		Ref:   engine.GetRef().String(),
		Input: engine,
	}

	if _, ok := r.URL.Query()["references"]; ok {
		body.OutputFormat = engine.GetDatasReferences()
		config.GetRefs()
	} else if config.References != nil {
		config = &core.Configs{
			Generic:  config.Generic,
			Specific: config.Specific,
		}
	}
	body.Config = config

	w.WriteJson(body)
}

// Set config imports
// PATCH /imports/:name/config?collection=COLLECTION_NAME
func (a *API) PatchImportConfig(w rest.ResponseWriter, r *rest.Request) {
	i := a.getImportByName(w, r)
	if i == nil {
		return
	}

	collection := a.getSingleCollectionByQuery(w, r)
	if collection == nil {
		return
	}

	var newConfigs core.Configs
	err := r.DecodeJsonPayload(&newConfigs)
	if err != nil {
		rest.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	fmt.Printf("[DECODE CONFIG] %+v\n", newConfigs)

	if err := i.SetConfig(collection, &newConfigs); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Handle import params
// PUT /imports/:name/params/:param
func (a *API) PutImportParams(w rest.ResponseWriter, r *rest.Request) {
	// Check if 'name' is an existing import
	name := r.PathParam("name")
	i, err := a.Classify.GetImportByName(name)
	if err != nil {

		// Otherwise try to create an ephemerous import with type specified
		i, err = core.NewImport(name)
		if err != nil {
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	var body json.RawMessage
	if err := r.DecodeJsonPayload(&body); err != nil {
		rest.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	res, err := i.GetParam(r.PathParam("param"), body)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteJson(res)
	return
}
