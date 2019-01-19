package core

import (
	"encoding/json"
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ohohleo/classify/collections"
	"golang.org/x/net/websocket"
	"net/http"
)

type ApiCollection struct {
	Name   string          `json:"name"`
	Ref    string          `json:"ref"`
	Config json.RawMessage `json:"config,omitempty" `
	Params json.RawMessage `json:"params,omitempty"`
}

// AddCollection adds new collection by API
// POST /collections
func (c *Classify) ApiPostCollection(w rest.ResponseWriter, r *rest.Request) {

	var body ApiCollection
	if err := r.DecodeJsonPayload(&body); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if body.Name == "" {
		rest.Error(w, "expected collection name",
			http.StatusBadRequest)
		return
	}

	if body.Ref == "" {
		rest.Error(w, "expected collection ref",
			http.StatusBadRequest)
		return
	}

	// Check the collection type
	ref, ok := collections.REF_STR2IDX[body.Ref]
	if ok == false {
		rest.Error(w, "invalid collection ref",
			http.StatusBadRequest)
		return
	}

	// Create new collection
	collection, err := c.AddCollection(body.Name, ref, body.Config, body.Params)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Store collection if enable
	if c.database != nil {

		err := collection.Store2DB(c.database)
		if err != nil {
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetCollections returns the name & the specificity of each collection
// GET /collections
func (c *Classify) ApiGetCollections(w rest.ResponseWriter, r *rest.Request) {

	result := make([]ApiCollection, len(c.collections))
	i := 0

	for name, c := range c.collections {

		result[i] = ApiCollection{
			Name: name,
			Ref:  c.engine.GetRef().String(),
		}

		i++
	}

	if err := w.WriteJson(result); err != nil {
		fmt.Printf("Get collections issue : %s\n", err.Error())
	}
}

func (c *Classify) getCollectionByName(w rest.ResponseWriter, r *rest.Request) *Collection {

	collection, err := c.GetCollection(r.PathParam("name"))
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}

	return collection
}

// GetCollectionByName returns the content of each collection
// GET /collections/:name
func (c *Classify) ApiGetCollectionByName(w rest.ResponseWriter, r *rest.Request) {

	if collection := c.getCollectionByName(w, r); collection != nil {
		w.WriteJson(collection)
	}
}

// PatchCollection modify the collection specified
// PATCH /collections/:name
func (c *Classify) ApiPatchCollection(w rest.ResponseWriter, r *rest.Request) {

	var body ApiCollection
	if err := r.DecodeJsonPayload(&body); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	isModified, err := c.ModifyCollection(r.PathParam("name"),
		body.Name, body.Ref)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if isModified {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusNotModified)
}

// DeleteCollectionByName delete the collection specified
// DELETE /collections/:name
func (c *Classify) ApiDeleteCollectionByName(w rest.ResponseWriter, r *rest.Request) {

	if err := c.DeleteCollection(r.PathParam("name")); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ApiGetCollectionConfig display actual configuration parameter
// GET /collections/:name/config[?references]
func (c *Classify) ApiGetCollectionConfig(w rest.ResponseWriter, r *rest.Request) {

	// Check the collection exist
	collection := c.getCollectionByName(w, r)
	if collection == nil {
		return
	}

	if collection.config == nil {
		rest.Error(w, "no collection config found",
			http.StatusBadRequest)
		return
	}

	// From the url query list
	values := r.URL.Query()

	_, ok := values["references"]
	w.WriteJson(collection.config.Get(ok))
}

// ApiPatchCollectionConfig mofify configuration parameters
// PATCH /collections/:name/config
func (c *Classify) ApiPatchCollectionConfig(w rest.ResponseWriter, r *rest.Request) {

	// Check the collection exist
	collection := c.getCollectionByName(w, r)
	if collection == nil {
		return
	}

	if err := r.DecodeJsonPayload(collection.config); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Store collection if enable
	if c.database != nil {

		if err := collection.StoreConfig2DB(c.database); err != nil {
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

// PUT /collections/:name/config/:param
func (c *Classify) ApiPutCollectionConfigParam(w rest.ResponseWriter, r *rest.Request) {

	// Check the collection exist
	collection := c.getCollectionByName(w, r)
	if collection == nil {
		return
	}

	path := r.PathParam("path")
	param := r.PathParam("param")

	var body json.RawMessage
	err := r.DecodeJsonPayload(&body)
	if err != nil {
		rest.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	res, err := collection.config.GetParam(path, param, body)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteJson(res)
	return
}

// GET /collections/:name/buffers
func (c *Classify) ApiGetCollectionBuffers(w rest.ResponseWriter, r *rest.Request) {

	// Check the collection exist
	collection := c.getCollectionByName(w, r)
	if collection == nil {
		return
	}

	w.WriteJson(collection.GetBuffer())
}

// DELETE /collections/:name/buffers
func (c *Classify) ApiDeleteCollectionBuffers(w rest.ResponseWriter, r *rest.Request) {

	// Check the collection exist
	collection := c.getCollectionByName(w, r)
	if collection == nil {
		return
	}

	collection.ResetBuffer()

	w.WriteHeader(http.StatusNoContent)
}

// GET /collections/:name/buffers/:id
func (c *Classify) ApiGetCollectionSingleBuffer(w rest.ResponseWriter, r *rest.Request) {

	// Check the collection exist
	collection := c.getCollectionByName(w, r)
	if collection == nil {
		return
	}
}

// PATCH /collections/:name/buffers/:id
func (c *Classify) ApiPatchCollectionSingleBuffer(w rest.ResponseWriter, r *rest.Request) {

	// Check the collection exist
	collection := c.getCollectionByName(w, r)
	if collection == nil {
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// POST /collections/:name/buffers/:id/validate
func (c *Classify) ApiValidateCollectionSingleBuffer(w rest.ResponseWriter, r *rest.Request) {

	// Check the collection exist
	collection := c.getCollectionByName(w, r)
	if collection == nil {
		return
	}

	err := collection.engine.Validate(r.PathParam("id"), json.NewDecoder(r.Body))
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Send next item if needed

	w.WriteHeader(http.StatusNoContent)
}

// DELETE /collections/:name/buffers/:id
func (c *Classify) ApiDeleteCollectionSingleBuffer(w rest.ResponseWriter, r *rest.Request) {

	// Check the collection exist
	collection := c.getCollectionByName(w, r)
	if collection == nil {
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GET /collections/:name/items
func (c *Classify) ApiGetCollectionItems(w rest.ResponseWriter, r *rest.Request) {

	// Check the collection exist
	collection := c.getCollectionByName(w, r)
	if collection == nil {
		return
	}

	w.WriteJson(collection.GetItems())
}

// DELETE /collections/:name/items
func (c *Classify) ApiDeleteCollectionItems(w rest.ResponseWriter, r *rest.Request) {

	// Check the collection exist
	collection := c.getCollectionByName(w, r)
	if collection == nil {
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// StartCollection launch the analysis of the collection
// PUT /collections/:name/start
func (c *Classify) ApiStartCollection(w rest.ResponseWriter, r *rest.Request) {

	// // Check the collection exist
	// collection := c.getCollectionByName(w, r)
	// if collection == nil {
	// 	return
	// }

	// channel, err := collection.Start()
	// if err != nil {
	// 	rest.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	// go func() {
	// 	for {
	// 		if item, ok := <-channel; ok {
	// 			c.Server.Send(collection.GetName(), "newFile", item)
	// 			continue
	// 		}
	// 		break
	// 	}
	// }()

	w.WriteHeader(http.StatusNoContent)
}

// StopCollection stop the analysis of the collection
// PUT /collections/:name/stop
func (c *Classify) ApiStopCollection(w rest.ResponseWriter, r *rest.Request) {

	w.WriteHeader(http.StatusNoContent)
}

func (c *Classify) ApiGetCollectionReferences(w rest.ResponseWriter, r *rest.Request) {

	collection := c.getCollectionByName(w, r)
	if collection == nil {
		return
	}

	w.WriteJson(References{
		Datas: collection.GetDatasReferences(),
	})
}

type Websocket interface {
	Handle(ws *websocket.Conn) error
}

type GetReferences struct {
	Websites    []string `json:"websites"`
	Collections []string `json:"collections"`
	Imports     []string `json:"imports"`
	Exports     []string `json:"exports"`
}

// GetReferences returns the website & type of collections available
// GET /references
func (c *Classify) ApiGetReferences(w rest.ResponseWriter, r *rest.Request) {

	w.WriteJson(GetReferences{
		Collections: c.GetCollectionRefs(),
		Imports:     c.GetImportRefs(),
		Exports:     c.GetExportRefs(),
		Websites:    c.GetWebsites(),
	})
}
