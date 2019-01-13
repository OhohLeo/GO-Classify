package core

import (
	"github.com/ohohleo/classify/reference"
)

type Configs struct {
	Collection *Collection                 `json:"-"`
	Generic    *GenericConfig              `json:"generic"`
	Specific   interface{}                 `json:"specific"`
	References map[string][]*reference.Ref `json:"references"`
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

	if c.Specific != nil {
		c.References["specific"] = reference.GetRefs(c.Specific)
	}
}

type GenericConfig struct {
	General struct {
		Enabled bool `json:"enabled"`
	} `json:"general"`

	Tweak *Tweak `json:"tweak"`
}

func NewGenericConfig() *GenericConfig {
	config := &GenericConfig{}

	// Enable import by default
	config.General.Enabled = true

	return config
}

func (g *GenericConfig) SetTweak(tweak *Tweak) {
	g.Tweak = tweak
}
