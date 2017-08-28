package core

import (
	"encoding/json"
	"github.com/ant0ine/go-json-rest/rest"
	"net/http"
	"strconv"
)

// List all config exports
// GET /exports/config
func (c *Classify) ApiGetExportsConfig(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson(c.config.Exports)
}

// getExportIdsAndCollections get from Url parameters exports and the collections
func (c *Classify) getExportIdsAndCollections(r *rest.Request) (exports map[uint64]*Export, collections map[string]*Collection, err error) {

	// From the url query list
	values := r.URL.Query()

	var ids []uint64
	for _, idStr := range values["id"] {

		var id int
		id, err = strconv.Atoi(idStr)
		if err != nil {
			return
		}

		ids = append(ids, uint64(id))
	}

	// Check and get the export list
	exports, err = c.GetExportsByIds(ids)
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
	Type        string          `json:"type"`
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

	i, err := c.AddExport(body.Type, body.Params, collections)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteJson(i)
}

// List all the exports selected by id or by collections
// GET /exports?id=EXPORT_ID&collection=COLLECTION_NAME
func (c *Classify) ApiGetExports(w rest.ResponseWriter, r *rest.Request) {

	ids, collections, err := c.getExportIdsAndCollections(r)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := c.GetExports(ids, collections)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteJson(res)
}

// ApiDeleteExport remove specified export selected by id and by the
// collections
// DELETE /exports?id=EXPORT_ID&collection=COLLECTION_NAME
func (c *Classify) ApiDeleteExport(w rest.ResponseWriter, r *rest.Request) {

	ids, collections, err := c.getExportIdsAndCollections(r)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := c.DeleteExports(ids, collections); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Stop the analysis of the collection export
// PUT /exports/stop?id=EXPORT_ID&collection=COLLECTION_NAME
func (c *Classify) ApiStopExport(w rest.ResponseWriter, r *rest.Request) {

	ids, collections, err := c.getExportIdsAndCollections(r)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := c.StopExports(ids, collections); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
