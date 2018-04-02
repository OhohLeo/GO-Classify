package core

import (
	"testing"

	"github.com/ant0ine/go-json-rest/rest/test"
	"github.com/stretchr/testify/assert"
)

func TestApiImport(t *testing.T) {

	assert := assert.New(t)

	// Create API
	api, err := new(Classify).GetApi(nil)
	assert.Nil(err)

	// Create new Import

	// Create new collection

	// Create tweak
	recorded := test.RunRequest(t, api.MakeHandler(),
		test.MakeSimpleRequest("GET", "http://localhost/imports", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
}
