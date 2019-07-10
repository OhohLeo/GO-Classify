package api

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ohohleo/classify/core"
)

type GetReferences struct {
	Collections map[string]core.References `json:"collections"`
	Imports     map[string]core.References `json:"imports"`
	Exports     map[string]core.References `json:"exports"`
	Websites    []string                   `json:"websites"`
}

// GetReferences returns the website & type of collections available
// GET /references
func (a *API) GetReferences(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson(GetReferences{
		Collections: a.Classify.GetCollectionReferences(),
		Imports:     a.Classify.GetImportRefs(),
		Exports:     a.Classify.GetExportRefs(),
		Websites:    a.Classify.GetWebsites(),
	})
}
