package src

import (
	"github.com/kelseyhightower/envconfig"
	"log"
)

type Configuration struct {
	Token string `envconfig:"BOT_TOKEN"`
}

var Config Configuration

// readConfigFromENV reads data from environment variables
func readConfigFromENV() (err error) {
	log.Println("Looking for ENV configuration")
	err = envconfig.Process("BOT", &Config)
	if err!= nil{
		return err
	}
	return
}