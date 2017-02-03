package collections

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfigClean(t *testing.T) {

	assert := assert.New(t)

	config := &Config{
		Separators: []string{"."},
		Banned:     []string{"HDTS"},
	}

	result, banned := config.clean("test.HDTS.2016")
	assert.Equal([]string{"test", "2016"}, result)
	assert.Equal([]string{"HDTS"}, banned)

	config = &Config{
		Separators: []string{},
		Banned:     []string{"HDTS"},
	}

	result, banned = config.clean("test.HDTS.2016")
	assert.Equal([]string{"test..2016"}, result)
	assert.Equal([]string{"HDTS"}, banned)
}
