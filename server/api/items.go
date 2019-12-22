package api

import (
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ohohleo/classify/core"
)

func (a *API) getItemById(w rest.ResponseWriter, r *rest.Request) *core.Item {

	collection := a.getCollectionByName(w, r)
	if collection == nil {
		return nil
	}

	item, err := collection.GetItemByString(r.PathParam("id"))
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}

	return item
}

// GET /collections/:name/items/:id?content="..."
func (a *API) GetCollectionSingleItem(w rest.ResponseWriter, r *rest.Request) {

	// Check the collection and item exist
	item := a.getItemById(w, r)
	if item == nil {
		return
	}

	// Ask for specific item data (image, ...)
	contentName := r.URL.Query().Get("content")

	if contentName != "" {

		contentPath := item.GetContent(contentName)
		if contentPath == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		http.ServeFile(w.(http.ResponseWriter), r.Request, contentPath)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// PATCH /collections/:name/items/:id
func (a *API) PatchCollectionSingleItem(w rest.ResponseWriter, r *rest.Request) {

	// Check the collection exist
	collection := a.getCollectionByName(w, r)
	if collection == nil {
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DELETE /collections/:name/items/:id
func (a *API) DeleteCollectionSingleItem(w rest.ResponseWriter, r *rest.Request) {

	// Check the collection exist
	collection := a.getCollectionByName(w, r)
	if collection == nil {
		return
	}

	id, err := core.GetIdFromString(r.PathParam("id"))
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = collection.RemoveItem(id); err != nil {
		rest.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
