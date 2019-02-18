package collections

import (
	"encoding/json"
)

const (
	MOVIES Ref = iota
	SIMPLE
)

type Ref uint64

func (t Ref) String() string {
	return REF_IDX2STR[t]
}

var REF_IDX2STR = []string{
	"movies",
	"simple",
}

var REF_STR2IDX = map[string]Ref{
	REF_IDX2STR[MOVIES]: MOVIES,
	REF_IDX2STR[SIMPLE]: SIMPLE,
}

type Collection interface {
	Check(json.RawMessage) error
	GetRef() Ref
	Validate(string, *json.Decoder) error
}

type Build struct {
	CheckConfig func(json.RawMessage) error
	ForceCreate func() Collection
	Create      func(json.RawMessage, json.RawMessage) (Collection, error)
}
