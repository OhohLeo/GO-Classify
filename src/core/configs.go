package core

import (
	"github.com/ohohleo/classify/reference"
)

type Configs struct {
	Collection *Collection                 `json:"-"`
	Generic    *GenericConfig              `json:"generic"`
	Tweak      *Tweak                      `json:"tweak"`
	Specific   interface{}                 `json:"specific,omitempty"`
	References map[string][]*reference.Ref `json:"references,omitempty"`
}

func NewConfigs(collection *Collection, specialConfig interface{}) *Configs {
	return &Configs{
		Collection: collection,
		Generic:    NewGenericConfig(),
		Specific:   specialConfig,
	}
}

func (c *Configs) GetRefs() {
	if c.References != nil {
		return
	}

	c.References = make(map[string][]*reference.Ref)
	c.References["generic"] = reference.GetRefs(c.Generic)
	c.References["tweak"] = []*reference.Ref{}

	if c.Specific != nil {
		c.References["specific"] = reference.GetRefs(c.Specific)
	}
}

type GenericConfig struct {
	Enabled bool `json:"enabled"`
}

func NewGenericConfig() *GenericConfig {
	config := &GenericConfig{}

	// Enable import by default
	config.Enabled = true

	return config
}
