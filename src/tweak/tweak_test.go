package tweak

import (
	"encoding/json"
	//"fmt"
	"testing"

	"github.com/ohohleo/classify/data"
	"github.com/stretchr/testify/assert"
)

var tweakJson = `{
  "source": {
    "file": {
      "name": {
        "regexp": "([a-z]+)([0-9]+)"
      },
      "path": {
        "regexp": "(\\d{4}-\\d{2}-\\d{2})"
      }
    }
  },
  "destination": {
    "item": {
      "name": {
        "value": "$2 $1:1"
      }
    }
  }
}`

func TestTweakJson(t *testing.T) {

	assert := assert.New(t)

	tweak, err := New([]byte(tweakJson))
	assert.Nil(err)

	jsonRes, err := json.MarshalIndent(tweak, "", "  ")
	assert.Nil(err)

	assert.Equal(tweakJson, string(jsonRes))
}

func TestTweak(t *testing.T) {

	assert := assert.New(t)

	tweak, err := New([]byte(tweakJson))
	assert.Nil(err)

	file := &data.File{
		Name: "abcd 1234",
		Path: "/path/to/test/(2017-01-02) test",
	}

	//var results map[string]Results
	_, err = tweak.Tweak(map[string]data.Data{
		"file": file,
	})
	assert.Nil(err)
}
