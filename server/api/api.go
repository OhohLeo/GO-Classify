package api

import (
	"github.com/ohohleo/classify/core"
)

type API struct {
	Classify *core.Classify
	Server   *Server
}

func NewAPI(classify *core.Classify, url string) (*API, error) {
	a := &API{
		Classify: classify,
	}

	var err error
	if a.Server, err = a.NewServer(url); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *API) Start() error {
	return a.Server.Start()
}

func (a *API) Stop() {
	if a.Server != nil {
		a.Server.Stop()
	}
}

func (a *API) SendEvent(event interface{}) {
	a.Server.SendEvent(event)
}
