package api

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ohohleo/classify/collections"
	"github.com/ohohleo/classify/core"
)

type Collection struct {
	Name   string          `json:"name"`
	Ref    string          `json:"ref"`
	Config json.RawMessage `json:"config,omitempty" `
	Params json.RawMessage `json:"params,omitempty"`
}

// AddCollection adds new collection by API
// POST /collections
func (a *API) PostCollection(w rest.ResponseWriter, r *rest.Request) {

	var body Collection
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
	_, err := a.Classify.CreateCollection(body.Name, ref, body.Config, body.Params)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetCollections returns the name & the specificity of each collection
// GET /collections
func (a *API) GetCollections(w rest.ResponseWriter, r *rest.Request) {

	result := make([]Collection, len(a.Classify.Collections))
	i := 0

	for name, collection := range a.Classify.Collections {

		result[i] = Collection{
			Name: name,
			Ref:  collection.Engine.GetRef().String(),
		}

		i++
	}

	if err := w.WriteJson(result); err != nil {
		fmt.Printf("Get collections issue : %s\n", err.Error())
	}
}

func (a *API) getCollectionByName(w rest.ResponseWriter, r *rest.Request) *core.Collection {

	collection, err := a.Classify.GetCollection(r.PathParam("name"))
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}

	return collection
}

// GetCollectionByName returns the content of each collection
// GET /collections/:name
func (a *API) GetCollectionByName(w rest.ResponseWriter, r *rest.Request) {

	if collection := a.getCollectionByName(w, r); collection != nil {
		w.WriteJson(collection)
	}
}

// PatchCollection modify the collection specified
// PATCH /collections/:name
func (a *API) PatchCollection(w rest.ResponseWriter, r *rest.Request) {

	var body Collection
	if err := r.DecodeJsonPayload(&body); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	isModified, err := a.Classify.ModifyCollection(r.PathParam("name"),
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
func (a *API) DeleteCollectionByName(w rest.ResponseWriter, r *rest.Request) {

	if err := a.Classify.DeleteCollection(r.PathParam("name")); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetCollectionConfig display actual configuration parameter
// GET /collections/:name/config[?references]
func (a *API) GetCollectionConfig(w rest.ResponseWriter, r *rest.Request) {

	// Check the collection exist
	collection := a.getCollectionByName(w, r)
	if collection == nil {
		return
	}

	if collection.Config == nil {
		rest.Error(w, "no collection config found",
			http.StatusBadRequest)
		return
	}

	// From the url query list
	values := r.URL.Query()

	_, ok := values["references"]
	w.WriteJson(collection.Config.Get(ok))
}

// PatchCollectionConfig mofify configuration parameters
// PATCH /collections/:name/config
func (a *API) PatchCollectionConfig(w rest.ResponseWriter, r *rest.Request) {

	// Check the collection exist
	collection := a.getCollectionByName(w, r)
	if collection == nil {
		return
	}

	if err := r.DecodeJsonPayload(collection.Config); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// // Store collection if enable
	// if a.Classify.database != nil {

	// 	if err := collection.StoreConfig2DB(a.Classify.database); err != nil {
	// 		rest.Error(w, err.Error(), http.StatusBadRequest)
	// 		return
	// 	}
	// }

	w.WriteHeader(http.StatusNoContent)
}

// PUT /collections/:name/config/:param
func (a *API) PutCollectionConfigParam(w rest.ResponseWriter, r *rest.Request) {

	// // Check the collection exist
	// collection := a.getCollectionByName(w, r)
	// if collection == nil {
	// 	return
	// }

	// path := r.PathParam("path")
	// param := r.PathParam("param")

	// var body json.RawMessage
	// err := r.DecodeJsonPayload(&body)
	// if err != nil {
	// 	rest.Error(w, "invalid json body", http.StatusBadRequest)
	// 	return
	// }

	// res, err := collection.Config.GetParam(path, param, body)
	// if err != nil {
	// 	rest.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	// w.WriteJson(res)
	return
}

// GET /collections/:name/buffers
func (a *API) GetCollectionBuffers(w rest.ResponseWriter, r *rest.Request) {

	// Check the collection exist
	collection := a.getCollectionByName(w, r)
	if collection == nil {
		return
	}

	w.WriteJson(collection.GetBuffer())
}

// DELETE /collections/:name/buffers
func (a *API) DeleteCollectionBuffers(w rest.ResponseWriter, r *rest.Request) {

	// Check the collection exist
	collection := a.getCollectionByName(w, r)
	if collection == nil {
		return
	}

	collection.ResetBuffer()

	w.WriteHeader(http.StatusNoContent)
}

// GET /collections/:name/buffers/:id
func (a *API) GetCollectionSingleBuffer(w rest.ResponseWriter, r *rest.Request) {

	// Check the collection exist
	collection := a.getCollectionByName(w, r)
	if collection == nil {
		return
	}
}

// PATCH /collections/:name/buffers/:id
func (a *API) PatchCollectionSingleBuffer(w rest.ResponseWriter, r *rest.Request) {

	// Check the collection exist
	collection := a.getCollectionByName(w, r)
	if collection == nil {
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// POST /collections/:name/buffers/:id/validate
func (a *API) ValidateCollectionSingleBuffer(w rest.ResponseWriter, r *rest.Request) {

	// Check the collection exist
	collection := a.getCollectionByName(w, r)
	if collection == nil {
		return
	}

	err := collection.Engine.Validate(r.PathParam("id"), json.NewDecoder(r.Body))
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Send next item if needed

	w.WriteHeader(http.StatusNoContent)
}

// DELETE /collections/:name/buffers/:id
func (a *API) DeleteCollectionSingleBuffer(w rest.ResponseWriter, r *rest.Request) {

	// Check the collection exist
	collection := a.getCollectionByName(w, r)
	if collection == nil {
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GET /collections/:name/items
func (a *API) GetCollectionItems(w rest.ResponseWriter, r *rest.Request) {

	// Check the collection exist
	collection := a.getCollectionByName(w, r)
	if collection == nil {
		return
	}

	w.WriteJson(collection.GetItems())
}

// DELETE /collections/:name/items
func (a *API) DeleteCollectionItems(w rest.ResponseWriter, r *rest.Request) {

	// Check the collection exist
	collection := a.getCollectionByName(w, r)
	if collection == nil {
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// StartCollection launch the analysis of the collection
// PUT /collections/:name/start
func (a *API) StartCollection(w rest.ResponseWriter, r *rest.Request) {

	// // Check the collection exist
	// collection := a.getCollectionByName(w, r)
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
	// 			a.Classify.Server.Send(collection.GetName(), "newFile", item)
	// 			continue
	// 		}
	// 		break
	// 	}
	// }()

	w.WriteHeader(http.StatusNoContent)
}

// StopCollection stop the analysis of the collection
// PUT /collections/:name/stop
func (a *API) StopCollection(w rest.ResponseWriter, r *rest.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) GetCollectionReferences(w rest.ResponseWriter, r *rest.Request) {
	collection := a.getCollectionByName(w, r)
	if collection == nil {
		return
	}

	w.WriteJson(collection.GetReferences())
}

type Websocket interface {
	Handle(ws *websocket.Conn) error
}
