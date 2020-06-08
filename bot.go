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
		if markup != nil {
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
		button1 := tgbotapi.KeyboardButton{Text: "Начать"}
		mrkup := tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{button1})
		return `Я бот. Я был создан для того, чтобы поведать тебе сказку о принцессе Анне и об одном парне, который создал меня, по имени Р… впрочем ты и так уже знаешь кто он. Что ж… Начнем? Shall we?`, mrkup
	case "Начать":
		m := `В тридевятом царстве, в тридесятом государстве, жила-была принцесса по имени Аня. Она была необычайно красивой, доброй, отзывчивой и умной девушкой. Из-за разных перипетий в её жизни она немного грустила. Все мы так делаем время от времени. Но сегодня… Сегодня у неё был повод радоваться, ведь сегодня у неё – День Рождения! Светлый праздник у всего королевства в этот день. Все птички напевают её любимые мелодии на протяжении всего дня! И был у неё знакомый парень по имени...`
		button1 := tgbotapi.KeyboardButton{Text: "Петр"}
		button2 := tgbotapi.KeyboardButton{Text: "Рома"}
		button3 := tgbotapi.KeyboardButton{Text: "Женя"}

		mrkup := tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{button1, button2, button3})
		return m, mrkup
	case "Петр", "Женя":
		return fmt.Sprintf("По имени %v. И он был парнем, который просто жил да был да как-то пересекался с принцессой Аней. Что на счет другого парня?", msg.Text), nil
	case "Рома":
		button1 := tgbotapi.KeyboardButton{Text: "Далее"}
		m := `По имени Рома. Рома очень хотел поздравить Аню с днем рождения и очень хотел вызвать на её прекрасном лице улыбку радость в её сердце. И он очень хотел произвести на неё впечатление.`
		mrkup := tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{button1})
		return m, mrkup
	case "Далее":
		button1 := tgbotapi.KeyboardButton{Text: "Да"}
		button2 := tgbotapi.KeyboardButton{Text: "Нет"}

		m := "Принцесса сидела в своей башне и слушала пение птиц, как вдруг одна птичка залатела прям к ней в окно, села на подоконник рядом с ней и уставилась на неё. Это был ручной птенчик Ромы. Она его уже раньше встречала. Принцесса заметила, что к птенца ноге была привязана бумага. Раскрыть бумагу?"
		mrkup := tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{button1, button2})
		return m, mrkup
	case "Да":
		return `Хороший выбор! Вот то, что было написано на записке: "С Днем рождения, милая принцесса! https://youtu.be/6H-InGqgKzo"`, nil
	case "Нет":
		return `Ты уверенна, что хочешь разбить сердце этому парню? Я, конечно, всего лишь набор нулей и еденичек, но даже я чувствую как ты ему дорога. А я не обишаюсь, я ведь машина! Уж поверь. https://i.ytimg.com/vi/zUOxSYOzECU/maxresdefault.jpg`, nil
	default:
		return "Такого выбора нет...", nil
	}
}
