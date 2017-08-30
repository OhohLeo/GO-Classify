package data

import ()

const (
	SIMPLE Ref = iota
	FILE
	MOVIE
)

type Ref uint64

func (t Ref) String() string {
	return REF_IDX2STR[t]
}

var REF_IDX2STR = []string{
	"simple",
	"file",
	"movie",
}

var REF_STR2IDX = map[string]Ref{
	REF_IDX2STR[SIMPLE]: SIMPLE,
	REF_IDX2STR[FILE]:   FILE,
	REF_IDX2STR[MOVIE]:  MOVIE,
}

type Data interface {
	GetName() string
	GetRef() Ref
}
