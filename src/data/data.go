package data

import (
	"encoding/json"
)

const (
	SIMPLE Ref = iota
	FILE
	MOVIE
	EMAIL
	ATTACHMENT
)

type Ref uint64

func (t Ref) String() string {
	return REF_IDX2STR[t]
}

var REF_IDX2STR = []string{
	"simple",
	"file",
	"movie",
	"email",
	"attachment",
}

var REF_STR2IDX = map[string]Ref{
	REF_IDX2STR[SIMPLE]:     SIMPLE,
	REF_IDX2STR[FILE]:       FILE,
	REF_IDX2STR[MOVIE]:      MOVIE,
	REF_IDX2STR[EMAIL]:      EMAIL,
	REF_IDX2STR[ATTACHMENT]: ATTACHMENT,
}

type Data interface {
	GetName() string
	GetRef() Ref
}

type HasDependencies interface {
	GetDependencies() []Data
}

type HasConfig interface {
	GetConfig() Config
}

type Config interface {
	UpdateConfig(json.RawMessage) error
}

type OnCollection interface {
	OnCollection(Config) error
}
