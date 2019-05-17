package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

type Bot struct {
	chatService Chat
	botAPI      *tgbotapi.BotAPI
	botConfig   tgbotapi.UpdateConfig
}

// NewBot returns new Bot service
func NewBot(b *tgbotapi.BotAPI, u tgbotapi.UpdateConfig, c Chat) Bot {
	return Bot{
		chatService: c,
		botAPI:    b,
		botConfig: u,
	}
}

func (b *Bot) listenMessages() error {
	updates, err := b.botAPI.GetUpdatesChan(b.botConfig)
	if err != nil {
		return err
	}

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		message := update.Message
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		b.MessageParser(message)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		b.botAPI.Send(msg)
	}
	return nil
}

func (b *Bot) MessageParser(msg *tgbotapi.Message) {
	switch msg.Text {
	case Start:
		b.chatService.AddChatToPull(msg.Chat.ID)
	}
}
