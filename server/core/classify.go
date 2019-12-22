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
	events      chan *Event
	imports     map[string]*Import
	exports     map[string]*Export
	Collections map[string]*Collection
	websites    map[string]websites.Website
}

type Config struct {
	Collections map[string]json.RawMessage `json:"collections"`
	Imports     map[string]json.RawMessage `json:"imports"`
	Exports     map[string]json.RawMessage `json:"exports"`

	API struct {
		URL string `json:"url,omitempty"`
	} `json:"api,omitempty"`

	// Database configuration
	DataBase database.Config `json:"database"`

	Websites map[string]map[string]string `json:"websites"`
}

// Application startup
func NewClassify(config *Config) (c *Classify, events chan *Event, err error) {
	c = new(Classify)

	log.Println("Reading configuration ...")

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

	// Store config file
	c.config = config

	// Init events output channel
	events = make(chan *Event)
	c.events = events

	// Specify that the application start
	go func() {
		c.SendEvent("start", "", "", "")
	}()

	return
}

func getRandomId() uint64 {
	return uint64(rand.Int63())
}

type Event struct {
	Event  string      `json:"event"`
	Status string      `json:"status"`
	Name   string      `json:"name"`
	Data   interface{} `json:"data"`
}

func (c *Classify) SendEvent(event string, status string, name string, data interface{}) {
	c.events <- &Event{
		Event:  event,
		Status: status,
		Name:   name,
		Data:   data,
	}
}
