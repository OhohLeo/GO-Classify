package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestItemClean(t *testing.T) {

	assert := assert.New(t)

	item := &Item{
		Name: "Lion.2016.720p.BRRip.x264.AAC-ETRG",
	}

	item.SetCleanedName([]string{"BRRip", "AAC-ETRG"},
		[]string{"."})

	assert.Equal(item.CleanedName, "Lion 2016 720p x264")
}
