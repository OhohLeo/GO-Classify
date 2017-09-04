package core

import (
	"github.com/ohohleo/classify/data"
)

type Item struct {
	Id     uint64    `json:"id"`
	Ref    string    `json:"ref"`
	Engine data.Data `json:"data"`
}

func (i *Item) SetData(input data.Data) {
	i.Ref = input.GetRef().String()
	i.Engine = input
}
