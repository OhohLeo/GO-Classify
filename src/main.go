package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/ohohleo/classify/api"
	"github.com/ohohleo/classify/core"
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

	classify, events, err := core.NewClassify(&config)
	if err != nil {
		log.Fatal(err.Error())
	}

	if config.API.URL != "" {
		api, err := api.NewAPI(classify, config.API.URL)
		if err != nil {
			log.Fatal(err.Error())
		}

		go func() {
			for event := range events {
				api.SendEvent(event)
			}
		}()

		if err = api.Start(); err != nil {
			log.Fatal(err.Error())
		}

		return
	}
}
