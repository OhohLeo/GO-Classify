package core

import (
	"encoding/json"
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ohohleo/classify/exports"
	"net/http"
)

// List all config exports
// GET /exports/config
func (c *Classify) ApiGetExportsConfig(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson(c.config.Exports)
}

// getExportNamesAndCollections get from Url parameters exports and the collections
func (c *Classify) getExportNamesAndCollections(r *rest.Request) (exports map[string]*Export, collections map[string]*Collection, err error) {

	// From the url query list
	values := r.URL.Query()

	var exportList []string
	for _, name := range values["name"] {
		exportList = append(exportList, name)
	}

	// Check and get the export list
	exports, err = c.GetExportsByNames(exportList)
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

type ApiAddExportsBody struct {
	Name        string          `json:"name"`
	Ref         string          `json:"ref"`
	Collections []string        `json:"collections"`
	Params      json.RawMessage `json:"params"`
}

// PostCollectionExport add a new export to the collection specified
// POST /exports
func (c *Classify) ApiAddExport(w rest.ResponseWriter, r *rest.Request) {

	// Get export parameters
	var body ApiAddExportsBody
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
	ref, ok := exports.REF_STR2IDX[body.Ref]
	if ok == false {
		rest.Error(w, fmt.Sprintf("invalid ref '%s'", body.Ref), http.StatusBadRequest)
		return
	}

	e, err := c.AddExport(body.Name, ref, body.Params, collections)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Handle DB storage
	if err := e.Store2DB(c.database); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteJson(e)
}

// List all the exports selected by id or by collections
// GET /exports?name=EXPORT_NAME&collection=COLLECTION_NAME
func (c *Classify) ApiGetExports(w rest.ResponseWriter, r *rest.Request) {

	names, collections, err := c.getExportNamesAndCollections(r)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := c.GetExports(names, collections)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteJson(res)
}

// ApiDeleteExport remove specified export selected by id and by the
// collections
// DELETE /exports?name=EXPORT_NAME&collection=COLLECTION_NAME
func (c *Classify) ApiDeleteExport(w rest.ResponseWriter, r *rest.Request) {

	names, collections, err := c.getExportNamesAndCollections(r)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := c.DeleteExports(names, collections); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Force exportation
// PUT /exports/force?name=EXPORT_NAME&collection=COLLECTION_NAME
func (c *Classify) ApiForceExport(w rest.ResponseWriter, r *rest.Request) {

	names, collections, err := c.getExportNamesAndCollections(r)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := c.ForceExports(names, collections); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Stop the analysis of the collection export
// PUT /exports/stop?name=EXPORT_NAME&collection=COLLECTION_NAME
func (c *Classify) ApiStopExport(w rest.ResponseWriter, r *rest.Request) {

	names, collections, err := c.getExportNamesAndCollections(r)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := c.StopExports(names, collections); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
