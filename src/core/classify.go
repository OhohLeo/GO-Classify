package core

import (
	log "github.com/Sirupsen/logrus"
	"math/rand"
)

type Classify struct {
	config      *Config
	Server      *Server
	imports     map[string]Import
	collections map[string]Collection
}

type Config struct {

	// Configuration du serveur
	Server ServerConfig `json:"server"`

	// Liste les configurations par type d'importation
	Imports map[string]map[string][]string `json:"imports"`

	// Liste les configurations par type d'exportation
	Exports map[string]map[string][]string `json:"exports"`
}

// Application startup
func Start(config Config) (c *Classify, err error) {

	c = new(Classify)

	log.Info("Config check...")

	err = c.CheckImportsConfig(config.Imports)
	if err != nil {
		return
	}

	log.Info("Config OK")

	log.SetLevel(log.DebugLevel)

	// TODO: Reload all collections saved

	log.Info("Start Classify")

	// Create server
	server, err := c.CreateServer(config.Server)
	if err != nil {
		return
	}

	// Store server
	c.Server = server

	// Store config file
	c.config = &config

	return
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
