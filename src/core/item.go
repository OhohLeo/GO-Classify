package core

import (
	"github.com/ohohleo/classify/data"
)

type Item struct {
	Id     uint64 `json:"id"`
	engine data.Data
}

func (i *Item) SetData(input data.Data) {
	i.engine = input
}
