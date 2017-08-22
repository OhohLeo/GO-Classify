package core

import (
	log "github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
	"github.com/ohohleo/classify/collections"
	"github.com/ohohleo/classify/database"
	"github.com/ohohleo/classify/requests"
	"github.com/ohohleo/classify/websites"
	"math/rand"
)

type Classify struct {
	config      *Config
	database    *sqlx.DB
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

	// Activate database if enabled
	if config.DataBase.Enable {

		log.Info("Starting Database")

		// Establish database connection
		if c.database, err = config.DataBase.Connect(); err != nil {
			return
		}

		// Create basic database tables
		if err = database.Create(c.database, c); err != nil {
			return
		}

		// Insert collection type references
		if err = database.InsertRef(
			c.database, &collections.DB_REFS,
			collections.TYPE_IDX2STR); err != nil {
			return
		}

		// Retreive all stored collections
		var storedCollections []collections.DB
		storedCollections, err = collections.GetDBCollections(c.database)
		if err != nil {
			return
		}

		for _, stored := range storedCollections {

			// Get collection params
			var params collections.Params
			params, err = stored.GetParams()
			if err != nil {
				return
			}

			// Create new collection
			_, err = c.AddCollection(
				stored.Name, collections.Type(stored.Type), params.Websites)
			if err != nil {
				return
			}
		}

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

// Application stop
func (c *Classify) Stop() {
	c.Server.Stop()
}

func (m *Classify) GetDBTables() []*database.Table {

	return []*database.Table{
		&collections.DB_LIST,
		&collections.DB_REFS,
	}
}

const nameLetters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func getRandomName() string {

	b := make([]byte, 16)
	for i := range b {
		b[i] = nameLetters[rand.Int63()%int64(len(nameLetters))]
	}

	return string(b)
}
