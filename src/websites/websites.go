package websites

import (
	"github.com/ohohleo/classify/data"
)

const (
	IMDB Ref = iota
	TMDB
)

type Ref int

func (t Ref) String() string {
	return REF_IDX2STR[t]
}

var REF_IDX2STR = []string{
	"IMDB",
	"TMDB",
}

type Website interface {
	GetRef() Ref
	SetConfig(map[string]string) bool
	Search(string) chan data.Data
}
