package main

import (
	"fmt"
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
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
		message := update.Message
		if message.Text != "" {
			log.Printf("[%s] %s", message.From.UserName, message.Text)
		}

		text, markup := b.AnnaMessageParser(message)
		msg := tgbotapi.NewMessage(message.Chat.ID, text)
		if markup != nil{
			msg.ReplyMarkup = markup
		}

		msg.ReplyToMessageID = message.MessageID

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
		return "Started listening the news! Wait for new post to arive! And remember... Roman Loves YOU! <3"
	default:
		return msg.Text
	}
}

func (b *Bot) AnnaMessageParser(msg *tgbotapi.Message) (string, interface{}) {
	switch msg.Text {
	case Start:
		button1 := tgbotapi.KeyboardButton{
			Text: "Получить",
		}
		button2 := tgbotapi.KeyboardButton{
			Text: "Отказаться",
		}
		mrkup := tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{button1, button2})
		return `Приветствую, принцесса :) Я всего лишь бот, но я был сделан для того, чтобы направить тебя на получение небольшого подарка, сделанного одним парнем, который хотел бы вызвать на твоем прекрасном лице улыбку  радость в твоем сердце. Напиши "Получить"" если хочешь получить подарок, либо /dl_Отказаться , если не желаешь его принять.`,
			mrkup
	case "Получить":
		return `Хороший выбор! С Днем рождения! https://youtu.be/6H-InGqgKzo`, nil
	case "Отказаться":
		return `Ты уверенна, что хочешь разбить сердце этому парню? Я, конечно, всего лишь набор нулей и еденичек, но даже я чувствую как ты ему дорога. А я не обишаюсь, я ведь машина! Уж поверь. https://i.ytimg.com/vi/zUOxSYOzECU/maxresdefault.jpg`, nil
	default:
		return "Такого выбора нет...", nil
	}
}
