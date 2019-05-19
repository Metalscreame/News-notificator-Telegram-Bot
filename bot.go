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

		text := b.MessageParser(message)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
		msg.ReplyToMessageID = update.Message.MessageID

		_, err := b.botAPI.Send(msg)
		if err != nil {
			log.Printf("ERROR: bot send failed %v", err)
		}
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

func (b *Bot) MessageParser(msg *tgbotapi.Message) string {
	switch msg.Text {
	case Start:
		b.chatService.AddChatToPull(msg.Chat.ID)
		return "Started listening news! Wait for new post!"
	default:
		return msg.Text
	}
}
