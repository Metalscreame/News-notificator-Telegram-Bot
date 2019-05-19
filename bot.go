package main

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
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
		botAPI:      b,
		botConfig:   u,
	}
}

func (b *Bot) listenMessages() {
	updates, err := b.botAPI.GetUpdatesChan(b.botConfig)
	if err != nil {
		log.Panic(err)
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
}

func (b *Bot) SendToChats(chats map[int64]struct{}, msg string) (err error) {
	for id := range chatMap {
		msg := tgbotapi.NewMessage(id, msg)
		m, err := b.botAPI.Send(msg)
		if err != nil {
			return fmt.Errorf("msg: %v, err: %v", m, err)
		}
	}
	return
}

func (b *Bot) MessageParser(msg *tgbotapi.Message) {
	switch msg.Text {
	case Start:
		b.chatService.AddChatToPull(msg.Chat.ID)
	}
}
