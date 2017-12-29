package core

import (
	"time"
)

type BlackList map[string]int64

func NewBlackList() BlackList {
	return make(map[string]int64)
}

func (b BlackList) Add(name string) {
	b[name] = time.Now().Unix()
}

func (b BlackList) Match(name string) bool {
	_, ok := b[name]
	if ok {
		delete(b, name)
	}

	return ok
}
