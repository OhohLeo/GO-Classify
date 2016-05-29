package main

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ohohleo/classify/imports/directory"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
)

type ApiCollectionBody struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// AddCollection adds new collection by API
// POST /collections
func ApiPostCollection(w rest.ResponseWriter, r *rest.Request) {

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

	_, err := classify.AddCollection(body.Name, body.Type)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type APICollection struct {
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

	collections := make([]APICollection, len(classify.collections))
	i := 0

	for name, c := range classify.collections {

		collectionType := c.GetType()

		image, ok := imagesByCollectionType[collectionType]
		if ok {
			image = "www/img/collections/" + image
		}

		collections[i] = APICollection{
			Name:  name,
			Type:  collectionType,
			Image: image,
		}

		i++
	}

	w.WriteJson(&collections)
}

func GetCollectionByName(w rest.ResponseWriter, r *rest.Request) Collection {

	collection, err := classify.GetCollection(r.PathParam("name"))
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}

	return collection
}

// GetCollectionByName returns the content of each collection
// GET /collection/:name
func ApiGetCollectionByName(w rest.ResponseWriter, r *rest.Request) {

	if collection := GetCollectionByName(w, r); collection != nil {
		w.WriteJson(collection)
	}
}

// ApiPatchCollection modify the collection specified
// PATCH /collection/:name
func ApiPatchCollection(w rest.ResponseWriter, r *rest.Request) {

	var body ApiCollectionBody
	if err := r.DecodeJsonPayload(&body); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	isModified, err := classify.ModifyCollection(r.PathParam("name"),
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
func ApiDeleteCollectionByName(w rest.ResponseWriter, r *rest.Request) {

	if err := classify.DeleteCollection(r.PathParam("name")); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ApiPostCollectionImport add a new import to the collection specified
// POST /collection/:name/imports
func ApiPostCollectionImport(w rest.ResponseWriter, r *rest.Request) {

	// Check the collection exist
	collection := GetCollectionByName(w, r)
	if collection == nil {
		return
	}

	// Get mapping
	var m Mapping
	err := r.DecodeJsonPayload(&m)
	if err != nil {
		return
	}

	// Create the new import
	name, i, err := classify.CreateNewImport(m)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = collection.AddImport(name, i)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ApiGetCollectionImport list all the imports used by the collection specified
// GET /collection/:name/imports
func ApiGetCollectionImports(w rest.ResponseWriter, r *rest.Request) {

	// Check the collection exist
	collection := GetCollectionByName(w, r)
	if collection == nil {
		return
	}

	w.WriteJson(collection.GetImports())
}

// ApiDeleteCollectionImport remove specified imports used by the collection specified
// DELETE /collection/:name/imports/:import
func ApiDeleteCollectionImport(w rest.ResponseWriter, r *rest.Request) {

	// Check the collection exist
	collection := GetCollectionByName(w, r)
	if collection == nil {
		return
	}

	if err := collection.DeleteImport(r.PathParam("import")); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ApiLaunchCollectionImport launch the reading of specified imports
// used by the collection specified
// PUT /collection/:name/imports/:import/launch
func ApiLaunchCollectionImport(w rest.ResponseWriter, r *rest.Request) {

	// Check the collection exist
	collection := GetCollectionByName(w, r)
	if collection == nil {
		return
	}

	channel, err := collection.LaunchImport(r.PathParam("import"))
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	go func() {
		for {
			if input, ok := <-channel; ok {
				//Send(ws, "newFile", input)

				log.Printf("API %+v\n", input)
				continue
			}
			break
		}
	}()

	w.WriteHeader(http.StatusNoContent)
}

type Websocket interface {
	Handle(ws *websocket.Conn) error
}

type NewDirectory struct {
	Path string `json:"path"`
}

type APIGetReferences struct {
	Websites []string `json:"websites"`
	Types    []string `json:"types"`
}

// GetReferences returns the website & type of collections available
// GET /references
func ApiGetReferences(w rest.ResponseWriter, r *rest.Request) {

	w.WriteJson(APIGetReferences{
		Websites: classify.GetWebsites(),
		Types:    classify.GetCollectionTypes(),
	})
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
