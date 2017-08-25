package imports

import (
	"encoding/json"
)

const (
	EMAIL Ref = iota
	DIRECTORY
)

type Ref int

func (t Ref) String() string {
	return REF_IDX2STR[t]
}

var REF_IDX2STR = []string{
	"email",
	"directory",
}

var REF_STR2IDX = map[string]Ref{
	REF_IDX2STR[EMAIL]:     EMAIL,
	REF_IDX2STR[DIRECTORY]: DIRECTORY,
}

type Import interface {
	Check(map[string][]string, []string) error
	Start() (chan Data, error)
	Stop()
	GetRef() Ref
	Eq(Import) bool
}

type Data interface {
	GetRef() string
	String() string
	GetUniqKey() string
}

type BuildImport struct {
	CheckConfig func(config map[string][]string) error
	Create      func(json.RawMessage, map[string][]string, []string) (Import, error)
}
