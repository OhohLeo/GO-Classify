package main

import (
	log "github.com/Sirupsen/logrus"
)

type Classify struct {
	Collections map[string]Collection
}

type Collection struct {
}

// Application startup
func Start() {
	log.SetLevel(log.DebugLevel)
	log.Info("Start Classify")
	ServerStart()

}

// Application stop
func Stop() {
}
