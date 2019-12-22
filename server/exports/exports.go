package exports

import (
	"encoding/json"
	"github.com/ohohleo/classify/data"
)

const (
	FILE Ref = iota
)

type Ref uint64

func (t Ref) String() string {
	return REF_IDX2STR[t]
}

var REF_IDX2STR = []string{
	"file",
}

var REF_STR2IDX = map[string]Ref{
	REF_IDX2STR[FILE]: FILE,
}

type Export interface {
	GetRef() Ref
	GetDatasReferences() []data.Data
	CheckConfig(json.RawMessage) error
	OnInput(data data.Data) error
	Stop() error
	Eq(Export) bool
}

type Build struct {
	CheckConfig func(json.RawMessage) error
	ForceCreate func() Export
	Create      func(json.RawMessage, json.RawMessage, []string) (Export, error)
}
