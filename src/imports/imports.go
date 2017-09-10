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
	Check(map[string][]string, []string) error
	Start() (chan data.Data, error)
	Stop() error
	GetRef() Ref
	Eq(Import) bool
}

type BuildImport struct {
	CheckConfig func(config map[string][]string) error
	Create      func(json.RawMessage, map[string][]string, []string) (Import, interface{}, error)
}
