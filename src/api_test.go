package main

import (
	"github.com/ohohleo/classify/requests"
	"github.com/stretchr/testify/assert"
	// "golang.org/x/net/websocket"
	//"log"
	"testing"
)

const URL = "http://127.0.0.1:3333"

func TestApi(t *testing.T) {
	assert := assert.New(t)

	requests.New(2, false)

	var classify *Classify

	go func() {
		checkGetReferences(assert)
		checkPostCollection(assert)
		checkGetCollections(assert)
		checkGetCollectionByName(assert)
		checkStartCollection(assert)
		checkStopCollection(assert)
		checkDeleteCollection(assert)
		classify.Stop()
	}()

	classify = Start()
}

func checkGetReferences(assert *assert.Assertions) {
	var rsp APIGetReferences
	c, err := requests.Send("GET", URL+"/references", nil, nil, &rsp)
	assert.Nil(err)

	result, ok := <-c
	assert.True(ok)
	assert.Equal(200, result.Status)

	assert.Equal(APIGetReferences{
		Websites: []string{"IMDB"},
		Types:    []string{"movies"},
	}, rsp)
}

func checkPostCollection(assert *assert.Assertions) {
	var rsp map[string]string

	// Failure : collection type doesn't exist
	c, err := requests.Send("POST", URL+"/collections",
		map[string]string{
			"Content-Type": "application/json",
		},
		ApiCollectionBody{
			Name: "test",
			Type: "error",
		}, &rsp)
	assert.Nil(err)

	result, ok := <-c
	assert.True(ok)
	assert.Equal(400, result.Status)

	assert.Equal(map[string]string{
		"Error": "invalid collection type 'error'",
	}, rsp)

	// Success : collection created
	c, err = requests.Send("POST", URL+"/collections",
		map[string]string{
			"Content-Type": "application/json",
		},
		ApiCollectionBody{
			Name: "test",
			Type: "movies",
		}, nil)
	assert.Nil(err)

	result, ok = <-c
	assert.True(ok)
	assert.Equal(204, result.Status)

	// Failure : collection already created
	c, err = requests.Send("POST", URL+"/collections",
		map[string]string{
			"Content-Type": "application/json",
		},
		ApiCollectionBody{
			Name: "test",
			Type: "movies",
		}, &rsp)
	assert.Nil(err)

	result, ok = <-c
	assert.True(ok)
	assert.Equal(400, result.Status)

	assert.Equal(map[string]string{
		"Error": "collection 'test' already exists",
	}, rsp)
}

func checkGetCollections(assert *assert.Assertions) {

	var rsp []APICollection

	// Success : get collections list
	c, err := requests.Send("GET", URL+"/collections",
		nil, nil, &rsp)
	assert.Nil(err)

	result, ok := <-c
	assert.True(ok)
	assert.Equal(200, result.Status)

	assert.Equal([]APICollection{
		APICollection{
			Name: "test",
			Type: "movies",
		},
	}, rsp)
}

func checkGetCollectionByName(assert *assert.Assertions) {

	var rsp Collection

	// Success : get collection 'test'
	c, err := requests.Send("GET", URL+"/collections/test",
		nil, nil, &rsp)
	assert.Nil(err)

	result, ok := <-c
	assert.True(ok)
	assert.Equal(200, result.Status)

	// TODO get result

	var rspError map[string]string

	// Failure : collection 'test' doesn't exist
	c, err = requests.Send("GET", URL+"/collections/error",
		nil, nil, &rspError)
	assert.Nil(err)

	result, ok = <-c
	assert.True(ok)
	assert.Equal(400, result.Status)

	assert.Equal(map[string]string{
		"Error": "collection 'error' not existing",
	}, rspError)
}

func checkStartCollection(assert *assert.Assertions) {

	// Success : delete specified collection
	c, err := requests.Send("PUT", URL+"/collections/test/start",
		nil, nil, nil)
	assert.Nil(err)

	result, ok := <-c
	assert.True(ok)
	assert.Equal(204, result.Status)
}

func checkStopCollection(assert *assert.Assertions) {

	// Success : delete specified collection
	c, err := requests.Send("PUT", URL+"/collections/test/stop",
		nil, nil, nil)
	assert.Nil(err)

	result, ok := <-c
	assert.True(ok)
	assert.Equal(204, result.Status)
}

func checkPatchCollection(assert *assert.Assertions) {

	// Success : patch collection 'test'
	c, err := requests.Send("GET", URL+"/collections/test",
		nil, nil, &rsp)
	assert.Nil(err)

	result, ok := <-c
	assert.True(ok)
	assert.Equal(200, result.Status)
}

func checkDeleteCollection(assert *assert.Assertions) {

	// Success : delete specified collection
	c, err := requests.Send("DELETE", URL+"/collections/test",
		nil, nil, nil)
	assert.Nil(err)

	result, ok := <-c
	assert.True(ok)
	assert.Equal(204, result.Status)

	// Failure : the collection doesn't exist
	var rsp map[string]string

	c, err = requests.Send("DELETE", URL+"/collections/test",
		nil, nil, &rsp)
	assert.Nil(err)

	result, ok = <-c
	assert.True(ok)
	assert.Equal(400, result.Status)

	assert.Equal(map[string]string{
		"Error": "collection 'test' not existing",
	}, rsp)
}
