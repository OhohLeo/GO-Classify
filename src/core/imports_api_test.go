package core

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ant0ine/go-json-rest/rest/test"
	"github.com/stretchr/testify/assert"
)

type RequestTest struct {
	Name              string
	Method            string
	Url               string
	Payload           interface{}
	ExpectedCode      int
	ExpectedRegexBody *regexp.Regexp
	ExpectedBody      string
}

func GenericTest(t *testing.T, api *rest.Api, check *RequestTest) bool {

	return t.Run(check.Name, func(t *testing.T) {
		recorded := test.RunRequest(t, api.MakeHandler(),
			test.MakeSimpleRequest(check.Method, check.Url, check.Payload))
		recorded.CodeIs(check.ExpectedCode)
		recorded.ContentTypeIsJson()

		if check.ExpectedRegexBody != nil {

			body, err := test.DecodedBody(recorded.Recorder)
			if err != nil {
				t.Errorf("Body '%s' expected, got error: '%s'", check.ExpectedRegexBody, err)
			}

			if check.ExpectedRegexBody.Match(body) == false {
				t.Errorf("Body '%s' expected, got '%s'", check.ExpectedRegexBody, body)
			}

		} else {
			recorded.BodyIs(check.ExpectedBody)
		}
	})
}

func TestApiImport(t *testing.T) {

	assert := assert.New(t)

	// Create API
	api, err := new(Classify).GetApi(nil)
	assert.Nil(err)

	tests := []*RequestTest{
		&RequestTest{
			Name:         "Post new import with empty payload",
			Method:       http.MethodPost,
			Url:          "http://localhost/imports",
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: `{
  "Error": "invalid json body: JSON payload is empty"
}`,
		},
		&RequestTest{
			Name:   "Post new import with invalid collection",
			Method: http.MethodPost,
			Url:    "http://localhost/imports",
			Payload: ApiAddImportsBody{
				Name:        "fail",
				Ref:         "fail",
				Collections: []string{"fail"},
			},
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: `{
  "Error": "collection 'fail' not existing"
}`,
		},
		// Add collection
		&RequestTest{
			Name:   "Post new collection",
			Method: http.MethodPost,
			Url:    "http://localhost/collections",
			Payload: ApiCollection{
				Name: "collection",
				Ref:  "simple",
			},
			ExpectedCode: http.StatusNoContent,
		},
		&RequestTest{
			Name:   "Post new import with invalid ref",
			Method: http.MethodPost,
			Url:    "http://localhost/imports",
			Payload: ApiAddImportsBody{
				Name:        "fail",
				Ref:         "fail",
				Collections: []string{"collection"},
			},
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: `{
  "Error": "invalid ref 'fail'"
}`,
		},
		&RequestTest{
			Name:   "Post new directory import",
			Method: http.MethodPost,
			Url:    "http://localhost/imports",
			Payload: ApiAddImportsBody{
				Name:        "directory",
				Ref:         "directory",
				Collections: []string{"collection"},
			},
			ExpectedCode: http.StatusOK,
			ExpectedRegexBody: regexp.MustCompile(`{
  "id": \d+,
  "name": "directory"
}`),
		},
		&RequestTest{
			Name:   "Post new imap import fail",
			Method: http.MethodPost,
			Url:    "http://localhost/imports",
			Payload: ApiAddImportsBody{
				Name:        "imap",
				Ref:         "imap",
				Collections: []string{"collection"},
			},
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: `{
  "Error": "import 'imap' connection: dial tcp :0: connect: connection refused"
}`,
		},
	}

	for _, check := range tests {
		if GenericTest(t, api, check) == false {
			t.Fail()
		}
	}
}

var addGenericImport = []*RequestTest{
	&RequestTest{
		Name:   "generic post new collection",
		Method: http.MethodPost,
		Url:    "http://localhost/collections",
		Payload: ApiCollection{
			Name: "collection",
			Ref:  "simple",
		},
		ExpectedCode: http.StatusNoContent,
	},
	&RequestTest{
		Name:   "generic post new import",
		Method: http.MethodPost,
		Url:    "http://localhost/imports",
		Payload: ApiAddImportsBody{
			Name:        "directory",
			Ref:         "directory",
			Collections: []string{"collection"},
		},
		ExpectedCode: http.StatusOK,
		ExpectedRegexBody: regexp.MustCompile(`{
  "id": \d+,
  "name": "directory"
}`),
	},
}

func TestApiGetImport(t *testing.T) {

	assert := assert.New(t)

	// Create API
	api, err := new(Classify).GetApi(nil)
	assert.Nil(err)

	tests := []*RequestTest{
		&RequestTest{
			Name:         "Get with no import",
			Method:       http.MethodGet,
			Url:          "http://localhost/imports",
			ExpectedCode: http.StatusOK,
			ExpectedBody: "{}",
		},
		addGenericImport[0],
		addGenericImport[1],
		&RequestTest{
			Name:         "Get with no import",
			Method:       http.MethodGet,
			Url:          "http://localhost/imports",
			ExpectedCode: http.StatusOK,
			ExpectedBody: `{
  "directory": {
    "directory": {
      "path": "",
      "is_recursive": false
    }
  }
}`,
		},
	}

	for _, check := range tests {
		if GenericTest(t, api, check) == false {
			t.Fail()
		}
	}
}

// 	&RequestTest{
// 			Name:         "Delete with invalid import",
// 			Method:       http.MethodDelete,
// 			Url:          "http://localhost/imports?name=fail",
// 			ExpectedCode: http.StatusBadRequest,
// 			ExpectedBody: `{
//   "Error": "import 'fail' not found"
// }`,
// 		},
// 		&RequestTest{
// 			Name:         "Put start with invalid import",
// 			Method:       http.MethodPut,
// 			Url:          "http://localhost/imports/start?name=fail",
// 			ExpectedCode: http.StatusBadRequest,
// 			ExpectedBody: `{
//   "Error": "import 'fail' not found"
// }`,
// 		},
// 		&RequestTest{
// 			Name:         "Put stop with invalid import",
// 			Method:       http.MethodPut,
// 			Url:          "http://localhost/imports/stop?name=fail",
// 			ExpectedCode: http.StatusBadRequest,
// 			ExpectedBody: `{
//   "Error": "import 'fail' not found"
// }`,
// 		},
// 		&RequestTest{
// 			Name:         "Put Tweak with invalid import",
// 			Method:       http.MethodPut,
// 			Url:          "http://localhost/imports/fail/tweak",
// 			ExpectedCode: http.StatusBadRequest,
// 			ExpectedBody: `{
//   "Error": "import 'fail' not found"
// }`,
// 		},
// 		&RequestTest{
// 			Name:         "Get tweak with invalid import",
// 			Method:       http.MethodGet,
// 			Url:          "http://localhost/imports/fail/tweak",
// 			ExpectedCode: http.StatusBadRequest,
// 			ExpectedBody: `{
//   "Error": "import 'fail' not found"
// }`,
// 		},
// 		&RequestTest{
// 			Name:         "Put param with invalid import",
// 			Method:       http.MethodPut,
// 			Url:          "http://localhost/imports/fail/param/type",
// 			ExpectedCode: http.StatusBadRequest,
// 			ExpectedBody: `{
//   "Error": "import 'fail' not found"
// }`,
// 		},
