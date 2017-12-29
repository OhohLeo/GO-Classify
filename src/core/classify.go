package core

import (
	"encoding/json"
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
	imports     map[string]*Import
	exports     map[string]*Export
	collections map[string]*Collection
	websites    map[string]websites.Website
}

type Config struct {
	Collections map[string]json.RawMessage `json:"collections"`
	Imports     map[string]json.RawMessage `json:"imports"`
	Exports     map[string]json.RawMessage `json:"exports"`

	// Server configuration
	Server ServerConfig `json:"server"`

	// Database configuration
	DataBase database.Config `json:"database"`

	Websites map[string]map[string]string `json:"websites"`
}

// Application startup
func Start(config *Config) (c *Classify, err error) {

	c = new(Classify)

	log.Println("Reading configuration ...")

	// Check configurations
	err = c.CheckCollectionsConfig(config.Collections)
	if err != nil {
		return
	}

	err = c.CheckImportsConfig(config.Imports)
	if err != nil {
		return
	}

	err = c.CheckExportsConfig(config.Exports)
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
