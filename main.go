package main

import (
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	err := readConfigFromENV()
	if err != nil {
		log.Fatalf("Could not load config from env, err: %v", err)
	}

	bot, u, err := setupBot()
	if err != nil {
		log.Panic(err)
	}
	chat := NewChat()

	botClient := NewBot(bot, u, chat)

	err = botClient.listenMessages()
	if err != nil {
		log.Panic(err)
	}
}

func setupBot() (*tgbotapi.BotAPI, tgbotapi.UpdateConfig, error) {
	bot, err := tgbotapi.NewBotAPI(Config.Token)
	if err != nil {
		return nil, tgbotapi.UpdateConfig{}, err
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	return bot, u, nil
}
