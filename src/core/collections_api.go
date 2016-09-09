package core

import (
	"github.com/ant0ine/go-json-rest/rest"
	"golang.org/x/net/websocket"
	"net/http"
)

type ApiCollectionBody struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// AddCollection adds new collection by API
// POST /collections
func (c *Classify) ApiPostCollection(w rest.ResponseWriter, r *rest.Request) {

	var body ApiCollectionBody
	if err := r.DecodeJsonPayload(&body); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if body.Name == "" {
		rest.Error(w, "expected collection name",
			http.StatusBadRequest)
		return
	}

	if body.Type == "" {
		rest.Error(w, "expected collection type",
			http.StatusBadRequest)
		return
	}

	_, err := c.AddCollection(body.Name, body.Type)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type ApiCollection struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// GetCollections returns the name & the specificity of each collection
// GET /collections
func (c *Classify) ApiGetCollections(w rest.ResponseWriter, r *rest.Request) {

	collections := make([]ApiCollection, len(c.collections))
	i := 0

	for name, c := range c.collections {

		collectionType := c.GetType()

		collections[i] = ApiCollection{
			Name: name,
			Type: collectionType,
		}

		i++
	}

	w.WriteJson(&collections)
}

func (c *Classify) getCollectionByName(w rest.ResponseWriter, r *rest.Request) Collection {

	collection, err := c.GetCollection(r.PathParam("name"))
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}

	return collection
}

// GetCollectionByName returns the content of each collection
// GET /collection/:name
func (c *Classify) ApiGetCollectionByName(w rest.ResponseWriter, r *rest.Request) {

	if collection := c.getCollectionByName(w, r); collection != nil {
		w.WriteJson(collection)
	}
}

// PatchCollection modify the collection specified
// PATCH /collection/:name
func (c *Classify) ApiPatchCollection(w rest.ResponseWriter, r *rest.Request) {

	var body ApiCollectionBody
	if err := r.DecodeJsonPayload(&body); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	isModified, err := c.ModifyCollection(r.PathParam("name"),
		body.Name, body.Type)
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
// DELETE /collection/:name
func (c *Classify) ApiDeleteCollectionByName(w rest.ResponseWriter, r *rest.Request) {

	if err := c.DeleteCollection(r.PathParam("name")); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// StartCollection launch the analysis of the collection
// PUT /collection/:name/start
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
// PUT /collection/:name/stop
func (c *Classify) ApiStopCollection(w rest.ResponseWriter, r *rest.Request) {

	w.WriteHeader(http.StatusNoContent)
}

type Websocket interface {
	Handle(ws *websocket.Conn) error
}

type GetReferences struct {
	Websites []string `json:"websites"`
	Types    []string `json:"types"`
}

// GetReferences returns the website & type of collections available
// GET /references
func (c *Classify) ApiGetReferences(w rest.ResponseWriter, r *rest.Request) {

	w.WriteJson(GetReferences{
		Websites: c.GetWebsites(),
		Types:    c.GetCollectionTypes(),
	})
}
