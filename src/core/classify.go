package core

import (
	"github.com/ohohleo/classify/database"
	"github.com/ohohleo/classify/requests"
	"github.com/ohohleo/classify/websites"
	"log"
	"math/rand"
)

type Classify struct {
	config      *Config
	database    *database.Database
	requests    *requests.RequestsPool
	Server      *Server
	imports     map[uint64]*Import
	exports     map[uint64]*Export
	collections map[string]*Collection
	websites    map[string]websites.Website
}

type Config struct {

	// Server configuration
	Server ServerConfig `json:"server"`

	// Database configuration
	DataBase database.Config `json:"database"`

	// Liste les configurations par type d'importation
	Imports  map[string]map[string][]string `json:"imports"`
	Websites map[string]map[string]string   `json:"websites"`

	// Liste les configurations par type d'exportation
	Exports map[string]map[string][]string `json:"exports"`
}

// Application startup
func Start(config *Config) (c *Classify, err error) {

	c = new(Classify)

	log.Println("Reading configuration ...")

	// Check configurations
	err = c.CheckImportsConfig(config.Imports)
	if err != nil {
		return
	}

	// Retreive classify stored data
	if err = c.StartDB(config); err != nil {
		return
	}

	log.Println("Starting Classify")

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
	c.config = config

	// Specify that the application start
	go func() {
		c.SendEvent("start", "", "", "")
	}()

	return
}

// Stop application
func (c *Classify) Stop() {
	c.Server.Stop()
}

func getRandomId() uint64 {
	return uint64(rand.Int63())
}
