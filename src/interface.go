package main

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ohohleo/classify/imports/directory"
	"golang.org/x/net/websocket"
	"net/http"
)

type APIPostCollectionReq struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// AddCollection adds new collection by API
// POST /collections
func ApiPostCollection(w rest.ResponseWriter, r *rest.Request) {

	var body APIPostCollectionReq
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

	_, err := classify.AddCollection(body.Name, body.Type)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type APICollections struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Image string `json:"image"`
}

var imagesByCollectionType = map[string]string{
	"movies": "movies.png",
}

// GetCollections returns the name & the specificity of each collection
// GET /collections
func ApiGetCollections(w rest.ResponseWriter, r *rest.Request) {

	collections := make([]APICollections, len(classify.collections))
	i := 0

	for name, c := range classify.collections {

		collectionType := c.GetType()

		image, ok := imagesByCollectionType[collectionType]
		if ok {
			image = "www/img/collections/" + image
		}

		collections[i] = APICollections{
			Name:  name,
			Type:  collectionType,
			Image: image,
		}

		i++
	}

	w.WriteJson(&collections)
}

// GetCollectionByName returns the content of each collection
// GET /collection/:name
func ApiGetCollectionByName(w rest.ResponseWriter, r *rest.Request) {

	collection, err := classify.GetCollection(r.PathParam("name"))
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteJson(collection)
}

// DeleteCollectionByName delete the collection specified
// DELETE /collection/:name
func ApiDeleteCollectionByName(w rest.ResponseWriter, r *rest.Request) {

	if err := classify.DeleteCollection(r.PathParam("name")); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type Websocket interface {
	Handle(ws *websocket.Conn) error
}

type NewDirectory struct {
	Path string `json:"path"`
}

// OnNewDirectory handle new directory input interface
func (rsp *NewDirectory) Handle(ws *websocket.Conn) error {

	d := &directory.Directory{
		Path:        rsp.Path,
		IsRecursive: true,
	}

	c, err := d.Launch()

	if err != nil {
		SendError(ws, err)
		return err
	}

	for {
		newFile, ok := <-c
		if ok == false {
			return nil
		}

		Send(ws, "newFile", newFile)
	}

	return nil
}
