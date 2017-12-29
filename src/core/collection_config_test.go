package core

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

func TestGetBannedList(t *testing.T) {

	assert := assert.New(t)

	config := &Config{
		Separators: []string{},
		Banned:     []string{"HDTS"},
	}

	assert.Equal([]string{"HDTS"}, config.GetBannedList())

	config = &Config{
		Separators: []string{"."},
		Banned:     []string{"HDTS", "toto.titi"},
	}

	assert.Equal([]string{"HDTS", "toto", "titi"}, config.GetBannedList())

	config = &Config{
		Separators: []string{".", "-"},
		Banned:     []string{"HDTS", "toto.titi-tutu"},
	}

	assert.Equal([]string{"HDTS", "toto", "titi", "tutu"}, config.GetBannedList())
}
