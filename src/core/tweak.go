package core

import (
	"github.com/ohohleo/classify/tweak"
)

type HasTweak interface {
	GetTweak(*Collection) *tweak.Tweak
	SetTweak(*Collection, *tweak.Tweak) error

	GetDatas() map[string]interface{}
}

func (c *Classify) GetTweak(t HasTweak, collection *Collection) *tweak.Tweak {

	res := t.GetTweak(collection)
	if res == nil {
		return new(tweak.Tweak)
	}

	return res
}

func (c *Classify) SetInputTweak(in HasTweak, collection *Collection, new *tweak.Tweak) (err error) {

	// Check tweak compatibility
	if err = new.Check(in.GetDatas(), collection.GetDatas()); err != nil {
		return
	}

	// Store tweak
	return in.SetTweak(collection, new)
}

func (c *Classify) SetExportTweak(out HasTweak, collection *Collection, new *tweak.Tweak) (err error) {

	// Check tweak compatibility
	if err = new.Check(collection.GetDatas(), out.GetDatas()); err != nil {
		return
	}

	// Store tweak
	return out.SetTweak(collection, new)
}
