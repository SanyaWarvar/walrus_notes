package main

import (
	"log"
	"wn/config"
	"wn/internal/app"
)

// @title           WALRUS NOTES API
// @version         1.0
// @description     This is walrus notes api service.

const configDir = "./config/main.yaml"

func main() {	
	cfg, err := config.NewConfig(configDir)

	if err != nil {
		log.Fatalf("Config error: %s", err)
	}
	// Run
	app.Run(cfg)
}
