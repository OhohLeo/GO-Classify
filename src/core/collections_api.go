package core

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ohohleo/classify/websites"
	"golang.org/x/net/websocket"
	"net/http"
)

type ApiCollection struct {
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Websites []string `json:"websites"`
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

	if body.Type == "" {
		rest.Error(w, "expected collection type",
			http.StatusBadRequest)
		return
	}

	// Check if the websites does existxsxs
	websites := make([]websites.Website, 0)
	for _, name := range body.Websites {

		website, err := c.AddWebsite(name)
		if err != nil {
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Add new website
		websites = append(websites, website)
	}

	// Create new collection
	collection, err := c.AddCollection(body.Name, body.Type)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Add websites to the collection created
	for _, website := range websites {
		collection.AddWebsite(website)
	}

	w.WriteHeader(http.StatusNoContent)
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
// DELETE /collections/:name
func (c *Classify) ApiDeleteCollectionByName(w rest.ResponseWriter, r *rest.Request) {

	if err := c.DeleteCollection(r.PathParam("name")); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ApiGetCollectionConfig display actual configuration parameter
// GET /collections/:name/config
func (c *Classify) ApiGetCollectionConfig(w rest.ResponseWriter, r *rest.Request) {

	// Check the collection exist
	collection := c.getCollectionByName(w, r)
	if collection == nil {
		return
	}

	w.WriteJson(collection.GetConfig())
}

type ApiCollectionConfig struct {
	Name   string   `json:"name"`
	Action string   `json:"action"`
	List   []string `json:"list"`
	Value  string   `json:"value"`
}

// ApiPatchCollectionConfig mofify configuration parameters
// PATCH /collections/:name/config
func (c *Classify) ApiPatchCollectionConfig(w rest.ResponseWriter, r *rest.Request) {

	// Check the collection exist
	collection := c.getCollectionByName(w, r)
	if collection == nil {
		return
	}

	var body ApiCollectionConfig
	if err := r.DecodeJsonPayload(&body); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if body.Action != "" {
		if err := collection.ModifyConfig(body.Name, body.Action, body.List); err != nil {
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else if err := collection.ModifyConfigValue(body.Name, body.Value); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
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

// PUT /collections/:name/buffers/:id/validate
func (c *Classify) ApiValidateCollectionSingleBuffer(w rest.ResponseWriter, r *rest.Request) {

	// Check the collection exist
	collection := c.getCollectionByName(w, r)
	if collection == nil {
		return
	}

	err := collection.ValidateBuffer(r.PathParam("id"))
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

// GET /collections/:name/items/:id
func (c *Classify) ApiGetCollectionSingleItem(w rest.ResponseWriter, r *rest.Request) {

	// Check the collection exist
	collection := c.getCollectionByName(w, r)
	if collection == nil {
		return
	}
}

// PATCH /collections/:name/items/:id
func (c *Classify) ApiPatchCollectionSingleItem(w rest.ResponseWriter, r *rest.Request) {

	// Check the collection exist
	collection := c.getCollectionByName(w, r)
	if collection == nil {
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DELETE /collections/:name/items/:id
func (c *Classify) ApiDeleteCollectionSingleItem(w rest.ResponseWriter, r *rest.Request) {

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
