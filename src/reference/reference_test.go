package reference

import (
	// "github.com/stretchr/testify/assert"
	"fmt"
	"testing"
)

type SimpleConfig struct {
	Bool       bool       `json:"bool"`
	Int        int        `json:"int"`
	String     string     `json:"string"`
	StringList StringList `json:"stringList"`
}

func TestSimpleConfig(t *testing.T) {

	simple := &SimpleConfig{
		Bool:       true,
		Int:        123,
		String:     "test",
		StringList: StringList{"toto", "titi", "tata"},
	}

	fmt.Printf("%+v", GetRefs(simple))

	// assert := assert.New(t)

}
