package tweak

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ohohleo/classify/core"
	"github.com/ohohleo/classify/data"
	"github.com/stretchr/testify/assert"
)

var TWEAK_JSON = `{
  "input": {
    "file": {
      "name": {
        "regexp": "([a-z]+)([0-9]+)"
      },
      "path": {
        "regexp": "(\\d{4}-\\d{2}-\\d{2})"
      }
    }
  },
  "output": {
    "item": {
      "name": {
        "value": ":file-name-0 :file-name-1 :file-path"
      }
    }
  }
}`

func TestTweakJson(t *testing.T) {

	assert := assert.New(t)

	tweak, err := New([]byte(TWEAK_JSON))
	assert.Nil(err)

	jsonRes, err := json.MarshalIndent(tweak, "", "  ")
	assert.Nil(err)

	assert.Equal(TWEAK_JSON, string(jsonRes))
}

func TestCheck(t *testing.T) {

	assert := assert.New(t)

	tweakJson := `
{
  "input": {
    "file": {
      "name": {
        "regexp": "([a-z]+)([0-9]+)"
      }
    }
  },
  "output": {
    "item": {
      "name": {
        "value": ":file-name-0 :file-name-1 :file-path"
      }
    }
  }
}`

	tweak, err := New([]byte(tweakJson))
	assert.Nil(err)

	type check struct {
		Input  map[string]interface{}
		Output map[string]interface{}
		Error  string
	}

	checks := []check{
		check{
			Input: map[string]interface{}{
				"file": &data.File{
					Name: "abcd1234",
					Path: "/path/// TODO: o/test/(2017-01-02) test",
				},
			},
			Output: map[string]interface{}{
				"item": &core.Item{
					Name: "test",
				},
			},
		},
		check{
			Input: map[string]interface{}{
				"file": &data.Email{
					Subject: "abcd1234",
				},
			},
			Output: map[string]interface{}{
				"item": &core.Item{
					Name: "test",
				},
			},
			Error: "invalid data file: field not found 'name'",
		},
	}

	for idx, check := range checks {

		err = tweak.Check(check.Input, check.Output)
		if err != nil && check.Error == "" {
			assert.Fail(fmt.Sprintf("error not expected but get error '%s' at %d", err.Error(), idx))
			return
		}

		if err == nil && check.Error != "" {
			assert.Fail(fmt.Sprintf("error expected '%s' but get no error at %d", check.Error, idx))
			return
		}

		if check.Error != "" {
			assert.Equal(check.Error, err.Error(), fmt.Sprintf("incompatible error at %d", idx))
			continue
		}
	}
}

func TestTweak(t *testing.T) {

	assert := assert.New(t)

	tweak, err := New([]byte(TWEAK_JSON))
	assert.Nil(err)

	type check struct {
		Input   map[string]interface{}
		Error   string
		Results map[string]map[string]string
	}

	checks := []check{
		check{
			Input: map[string]interface{}{
				"file": &data.File{
					Name: "abcd1234",
					Path: "/path/// TODO: o/test/(2017-01-02) test",
				},
			},
			Results: map[string]map[string]string{
				"item": map[string]string{
					"name": "abcd 1234 2017-01-02",
				},
			},
		},
	}

	var results map[string]map[string]string
	for idx, check := range checks {

		results, err = tweak.Tweak(check.Input)
		if err != nil && check.Error == "" {
			assert.Fail(fmt.Sprintf("error not expected but get error '%s' at %d", err.Error(), idx))
			return
		}

		if err == nil && check.Error != "" {
			assert.Fail(fmt.Sprintf("error expected '%s' but get no error at %d", check.Error, idx))
			return
		}

		if check.Error != "" {
			assert.Equal(check.Error, err.Error(), fmt.Sprintf("incompatible error at %d", idx))
			continue
		}

		assert.Nil(err, fmt.Sprintf("no error expected at %d", idx))
		assert.Equal(check.Results, results, fmt.Sprintf("incompatible result at %d", idx))
	}

}
