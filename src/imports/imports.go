package imports

import (
	"encoding/json"
	"github.com/ohohleo/classify/data"
)

const (
	IMAP Ref = iota
	DIRECTORY
)

type Ref uint64

func (t Ref) String() string {
	return REF_IDX2STR[t]
}

var REF_IDX2STR = []string{
	"imap",
	"directory",
}

var REF_STR2IDX = map[string]Ref{
	REF_IDX2STR[IMAP]:      IMAP,
	REF_IDX2STR[DIRECTORY]: DIRECTORY,
}

type Import interface {
	CheckConfig(json.RawMessage) error
	Start() (chan data.Data, error)
	Stop() error
	GetRef() Ref
	Eq(Import) bool
}

type Build struct {
	CheckConfig func(json.RawMessage) error
	Create      func(json.RawMessage, json.RawMessage, []string) (Import, interface{}, error)
	GetParam    func(string, json.RawMessage) (interface{}, error)
}
