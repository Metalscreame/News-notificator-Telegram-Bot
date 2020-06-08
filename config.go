package main

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Configuration struct {
	Token        string `envconfig:"BOT_TOKEN"`
	AnnaBotToken string `envconfig:"BOT_ANNA_TOKEN"`
}

var Config Configuration

// readConfigFromENV reads data from environment variables
func readConfigFromENV() (err error) {
	log.Println("Looking for ENV configuration")
	err = envconfig.Process("BOT", &Config)
	if err != nil {
		return err
	}
	return
}
