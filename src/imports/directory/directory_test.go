package directory

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestReadDirectory(t *testing.T) {

	assert := assert.New(t)

	directory := &Directory{
		Path:        `/home/ohohleo`,
		IsRecursive: false,
	}

	c, err := directory.Launch()
	assert.Nil(err)

	for {

		if f, ok := <-c; ok {
			log.Printf("Received %s\n", f.FullPath)
			continue
		}

		break
	}
}
