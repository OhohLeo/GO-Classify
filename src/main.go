package main

import (
	"encoding/json"
	"flag"
	"github.com/ohohleo/classify/core"
	"log"
	"os"
)

func main() {

	// Get application flags
	var configPath, serverUrl string
	flag.StringVar(&configPath, "config", "config.json", "Config file path")
	flag.StringVar(&serverUrl, "server", "", "Server url")
	flag.Parse()

	var config core.Config

	// Check if config path does exists
	file, err := os.Open(configPath)
	if err == nil {

		// Decode config file
		decoder := json.NewDecoder(file)
		err := decoder.Decode(&config)
		if err != nil {
			log.Fatal("Error in config file '" +
				configPath + "': " + err.Error())
		}
	}

	// Priority for flags parameters
	if serverUrl != "" {
		config.Server.Url = serverUrl
	}

	// Start application
	classify, err := core.Start(config)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Start server
	classify.Server.Start()
}
