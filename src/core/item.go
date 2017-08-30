package core

import (
	"github.com/ohohleo/classify/data"
)

type Item struct {
	Id     uint64 `json:"id"`
	Ref    data.Ref
	engine data.Data
}
