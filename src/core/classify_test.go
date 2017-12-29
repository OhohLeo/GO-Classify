package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCollection(t *testing.T) {
	c := new(Classify)
	assert := assert.New(t)

	collection, err := c.AddCollection("test", "movies")
	if assert.Nil(err) {
		assert.NotNil(collection)
	}

	collection, err = c.GetCollection("test")
	if assert.Nil(err) {
		assert.NotNil(collection)
	}

	err = c.DeleteCollection("test")
	assert.Nil(err)

	assert.Equal([]string{
		"movies",
	}, c.GetCollectionTypes())

	assert.Equal([]string{
		"IMDB",
	}, c.GetWebsites())
}
