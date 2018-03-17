package core

import (
	"github.com/ant0ine/go-json-rest/rest"
	"net/http"
)

func (c *Classify) getItemById(w rest.ResponseWriter, r *rest.Request) *Item {

	collection := c.getCollectionByName(w, r)
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
func (c *Classify) ApiGetCollectionSingleItem(w rest.ResponseWriter, r *rest.Request) {

	// Check the collection and item exist
	item := c.getItemById(w, r)
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

	if err := collection.RemoveItem(r.PathParam("id")); err != nil {
		rest.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
