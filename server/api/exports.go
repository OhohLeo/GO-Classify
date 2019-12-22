package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ohohleo/classify/core"
	"github.com/ohohleo/classify/exports"
)

// getExportByName get from Url parameters export
func (a *API) getExportByName(w rest.ResponseWriter, r *rest.Request) *core.Export {

	i, err := a.Classify.GetExportByName(r.PathParam("name"))
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}

	return i
}

// getExportNamesAndCollections get from Url parameters exports and the collections
func (a *API) getExportNamesAndCollections(r *rest.Request) (exports map[string]*core.Export, collections map[string]*core.Collection, err error) {

	// From the url query list
	values := r.URL.Query()

	var exportList []string
	for _, name := range values["name"] {
		exportList = append(exportList, name)
	}

	// Check and get the export list
	exports, err = a.Classify.GetExportsByNames(exportList)
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

type AddExportsBody struct {
	Name        string          `json:"name"`
	Ref         string          `json:"ref"`
	Collections []string        `json:"collections"`
	Params      json.RawMessage `json:"params"`
}

// PostCollectionExport add a new export to the collection specified
// POST /exports
func (a *API) AddExport(w rest.ResponseWriter, r *rest.Request) {

	// Get export parameters
	var body AddExportsBody
	err := r.DecodeJsonPayload(&body)
	if err != nil {
		rest.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	// Check and get the collection list
	collections, err := a.Classify.GetCollectionsByNames(body.Collections)
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

	e, err := a.Classify.CreateExport(body.Name, ref, body.Params, collections)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteJson(e)
}

// List all the exports selected by id or by collections
// GET /exports?name=EXPORT_NAME&collection=COLLECTION_NAME
func (a *API) GetExports(w rest.ResponseWriter, r *rest.Request) {

	names, collections, err := a.getExportNamesAndCollections(r)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := a.Classify.GetExports(names, collections)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteJson(res)
}

// DeleteExport remove specified export selected by id and by the
// collections
// DELETE /exports?name=EXPORT_NAME&collection=COLLECTION_NAME
func (a *API) DeleteExport(w rest.ResponseWriter, r *rest.Request) {

	names, collections, err := a.getExportNamesAndCollections(r)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := a.Classify.DeleteExports(names, collections); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Force exportation
// PUT /exports/force?name=EXPORT_NAME&collection=COLLECTION_NAME
func (a *API) ForceExport(w rest.ResponseWriter, r *rest.Request) {

	names, collections, err := a.getExportNamesAndCollections(r)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := a.Classify.ForceExports(names, collections); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Stop the analysis of the collection export
// PUT /exports/stop?name=EXPORT_NAME&collection=COLLECTION_NAME
func (a *API) StopExport(w rest.ResponseWriter, r *rest.Request) {

	names, collections, err := a.getExportNamesAndCollections(r)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := a.Classify.StopExports(names, collections); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// List all config exports
// GET /exports/:name/config
func (a *API) GetExportConfig(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson(nil)
}

// Set config exports
// PATCH /exports/:name/config
func (a *API) PatchExportConfig(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson(nil)
}

// Handle export params
// PUT /exports/:name/params/:params
func (a *API) PutExportParams(w rest.ResponseWriter, r *rest.Request) {

	// param := r.PathParam("param")
	// name := r.PathParam("name")

	// // Récupération du type de l'exportation
	// i, err := a.Classify.GetExportByName(name)
	// if err == nil {
	// 	name = i.engine.GetRef().String()
	// }

	// newExport, ok := newExports[name]
	// if ok == false {
	// 	rest.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	// var body json.RawMessage
	// err = r.DecodeJsonPayload(&body)
	// if err != nil {
	// 	rest.Error(w, "invalid json body", http.StatusBadRequest)
	// 	return
	// }

	// res, err := newExport.GetParam(param, body)
	// if err != nil {
	// 	rest.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	// w.WriteJson(res)
	return
}
