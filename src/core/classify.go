package core

import (
	log "github.com/Sirupsen/logrus"
	"math/rand"
)

type Classify struct {
	Server      *Server
	imports     map[string]Import
	collections map[string]Collection
}

// Application startup
func Start() (*Classify, error) {

	c := new(Classify)

	// TODO: Reload all collections saved

	log.SetLevel(log.DebugLevel)

	log.Info("Start Classify")

	server, err := c.CreateServer()
	if err != nil {
		return nil, err
	}

	// Store server
	c.Server = server

	return c, nil
}

// Application stop
func (c *Classify) Stop() {
	c.Server.Stop()
}

const nameLetters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func getRandomName() string {

	b := make([]byte, 16)
	for i := range b {
		b[i] = nameLetters[rand.Int63()%int64(len(nameLetters))]
	}

	return string(b)
}
