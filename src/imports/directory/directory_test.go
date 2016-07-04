package directory

import (
	"github.com/ohohleo/classify/imports"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestReadDirectory(t *testing.T) {

	assert := assert.New(t)

	directory := &Directory{
		Path:        `/home/lmartin`,
		IsRecursive: false,
	}

	c, err := directory.Start()
	if assert.Nil(err) {
		for {

			if data, ok := <-c; ok {
				if f, ok := data.(imports.File); ok {
					log.Printf("Received %s\n", f.FullPath)
					continue

				}
			}

			break
		}
	}
}
