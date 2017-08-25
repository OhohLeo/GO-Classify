package core

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ohohleo/classify/collections"
	"github.com/ohohleo/classify/database"
	"github.com/ohohleo/classify/imports"
	"github.com/ohohleo/classify/requests"
	"github.com/ohohleo/classify/websites"
	"math/rand"
)

type Classify struct {
	config      *Config
	database    *database.Database
	requests    *requests.RequestsPool
	Server      *Server
	imports     map[string]Import
	exports     map[string]Export
	collections map[string]Collection
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

	log.Info("Reading configuration ...")

	err = c.CheckImportsConfig(config.Imports)
	if err != nil {
		return
	}

	log.SetLevel(log.DebugLevel)

	if err = c.StartDB(config); err != nil {
		return
	}

	// TODO: Read all imports

	log.Info("Starting Classify")

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

// StartDB init db and retreive all stored data
func (c *Classify) StartDB(config *Config) (err error) {

	// Activate database if enabled
	c.database, err = database.New(config.DataBase)
	if c.database == nil {
		log.Info("Database disable")
		return
	}

	log.Info("Starting Database")

	// Init database tables
	if err = collections.INIT_DB(c.database); err != nil {
		return
	}

	if err = imports.INIT_DB(c.database); err != nil {
		return
	}

	// Create tables
	if err = c.database.Create(); err != nil {
		return
	}

	// Insert all references
	if err = collections.INIT_REF_DB(c.database); err != nil {
		return
	}

	if err = imports.INIT_REF_DB(c.database); err != nil {
		return
	}

	// Retreive all stored collections
	err = collections.RetreiveDBCollections(c.database,
		func(id uint64, name string, ref collections.Ref, params collections.Params) (err error) {
			collection, err := c.AddCollection(name, ref, params.Websites)
			if err != nil {
				collection.SetId(id)
			}
			return
		})
	if err != nil {
		return
	}

	return
}

const nameLetters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func getRandomName() string {

	b := make([]byte, 16)
	for i := range b {
		b[i] = nameLetters[rand.Int63()%int64(len(nameLetters))]
	}

	return string(b)
}
