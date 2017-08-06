package core

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ohohleo/classify/requests"
	"github.com/ohohleo/classify/websites"
	"math/rand"
)

type Classify struct {
	config      *Config
	requests    *requests.RequestsPool
	Server      *Server
	imports     map[string]Import
	exports     map[string]Export
	collections map[string]Collection
	websites    map[string]websites.Website
}

type Config struct {

	// Configuration du serveur
	Server ServerConfig `json:"server"`

	// Liste les configurations par type d'importation
	Imports map[string]map[string][]string `json:"imports"`

	Websites map[string]map[string]string `json:"websites"`

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

	log.SetLevel(log.DebugLevel)

	// TODO: Reload all collections saved

	// TODO: Reload all imports

	log.Info("Start Classify")

	// HTTP requests
	c.requests = requests.New(2, true)

	// Create server
	server, err := c.CreateServer(config.Server)
	if err != nil {
		return
	}

	// Store server
	c.Server = server

	// Store config file
	c.config = &config

	// Specify that the application start
	go func() {
		c.SendEvent("start", "", "", "")
	}()

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
