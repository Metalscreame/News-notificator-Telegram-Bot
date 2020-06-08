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
		return `Я бот. Я был создан для того, чтобы поведать тебе игру-сказку о принцессе Анне и об одном парне, который создал меня. Иногда он любит писать код и играть на гитаре, А зовут его Р… .Впрочем ты и так уже, наверное, знаешь кто он. Что ж… Начнем? Shall we? Не бойся выбирать разные варианты ответов`, mrkup
	case "Начать":
		m := `В тридевятом царстве, в тридесятом государстве, жила-была принцесса по имени Аня. Она была необычайно красивой, доброй, отзывчивой и умной девушкой. Из-за разных перипетий в её жизни она немного грустила. Все мы так делаем время от времени. Но сегодня… Сегодня у неё был повод радоваться, ведь сегодня у неё – День Рождения! Светлый праздник у всего королевства в этот день. Все птички напевают её любимые мелодии на протяжении всего дня! И был у неё знакомый парень по имени...`
		button1 := tgbotapi.KeyboardButton{Text: "Петр"}
		button2 := tgbotapi.KeyboardButton{Text: "Рома"}
		button3 := tgbotapi.KeyboardButton{Text: "Женя"}

		mrkup := tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{button1, button2, button3})
		return m, mrkup
	case "Петр", "Женя":
		return fmt.Sprintf("По имени %v. И он был парнем, который просто жил да был да как-то пересекался с принцессой Аней. Все же не наш Герой. Что на счет другого парня?", msg.Text), nil
	case "Рома":
		button1 := tgbotapi.KeyboardButton{Text: "Далее"}
		m := `По имени Рома. Рома очень хотел поздравить принцессу Аню с днем рождения и очень хотел вызвать на её прекрасном лице улыбку радость в сердце.Он хотел произвести на неё впечатление так как она ему очень нравилась. И он задумал сделать кое-что для принцессы, ибо он как-то услышал от неё то, что могло бы ей понравится, когда они общались вместе на берегу одной известной реки...`
		mrkup := tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{button1})
		return m, mrkup

	case "Далее":
		button1 := tgbotapi.KeyboardButton{Text: "Выйти на улицу"}
		button2 := tgbotapi.KeyboardButton{Text: "Подняться в башню"}
		mrkup := tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{button1, button2})
		m := "Вернемся к принцессе. Она подумывала чем бы занять себя в этот прекрасный день. На улице было очень тепло. Она хотела насладится природой либо снаружи, либо из высшей точки замка."
		return m, mrkup
	case "Выйти на улицу":
		button1 := tgbotapi.KeyboardButton{Text: "Подойти к пионам поближе"}
		button2 := tgbotapi.KeyboardButton{Text: "Ничего не делать"}
		m := "Принцесса вышла из своего замка и направилась в сад. В саду её встретил целый букет недавно распустившихся поздних цветов. Её взгляд упал на Пионы. "
		mrkup := tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{button1, button2})
		return m, mrkup
	case "Подойти к пионам поближе":
		button2 := tgbotapi.KeyboardButton{Text: "Подняться в башню"}
		m := "Её приманил к себе этот чудесный запах, который можно услышать только в майский день. Это запах свежых пионов. В её саду росли разные цветы, но именно в этот весенне-летний период её манило к ним. Она подобралась к ним поближе и сделала глубокий вдох. Её тело захватил этот чарующий запах. Она долго не могла отойти от этих прекрасных пионов. Спустя какое-то время она подняла голову и посмотрела на Башню."
		mrkup := tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{button2})
		return m, mrkup
	case "Ничего не делать":
		button2 := tgbotapi.KeyboardButton{Text: "Подняться в башню"}
		m := "Принцесса не стала любоваться пионами. Она решила посмотреть на свое королевство с высоты птичьего полета."
		mrkup := tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{button2})
		return m, mrkup
	case "Подняться в башню":
		button1 := tgbotapi.KeyboardButton{Text: "Да"}
		button2 := tgbotapi.KeyboardButton{Text: "Нет"}

		m := "Принцесса направилась в башню. Башня была так высока, что люди могли только поражаться этому чуду инженерной техники. Она очень любила там любоваться природой. Она подошла к окну и стала слушать пение птиц. Они пели так, словно они знали какой сегодня день. Это очень понравилось принцессе. Их песни были наполнены нотками радости, счастья и любви. Это согревало душу принцессе. Вдруг одна птичка залатела прям к ней в окно, села на подоконник рядом с ней и уставилась на неё. Это был ручной птенчик Ромы. Она его уже раньше встречала. Принцесса заметила, что у птенца на ноге была привязана бумага. Раскрыть бумагу?"
		mrkup := tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{button1, button2})
		return m, mrkup
	case "Да":
		m := `В бумаге она увидела послание. Она сразу узнала почерк. Это был почерк Ромы. Она стала читать послание. Там было написано "Иди к холму, что виден из этой башни на севере. В этом холме ты найдешь вход в пещеру. В пещере тебя будет ждать подарок, который делал от чистого сердца, чтобы поздравить тебя в этот день. Шибас (так зовут птенчика) проведет тебя ко входу"`
		button1 := tgbotapi.KeyboardButton{Text: "Отравиться в путь, который указан в послании"}
		button2 := tgbotapi.KeyboardButton{Text: "Продолжать смотреть в окно"}
		mrkup := tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{button1, button2})
		return m, mrkup
	case "Продолжать смотреть в окно":
		m := `Принцесса продолжила смотреть в окно. Но птичка все не улетала. Она будто ждала чего-то. Прекрасная Принцесса Аня не выдержала этого прелестного взгяла от Шибаса "Боже, какой он милый" подумала она. И она сказала "Кончено же я пойду". Шибас радостно чирикнул. Принцесса начала собираться и вышла на улицу.`
		button2 := tgbotapi.KeyboardButton{Text: "Идти к холму"}
		mrkup := tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{button2})
		return m, mrkup
	case "Идти к холму", "Отравиться в путь, который указан в послании":
		button2 := tgbotapi.KeyboardButton{Text: "Войти в пещеру"}
		button1 := tgbotapi.KeyboardButton{Text: "Не входить"}
		m := "Принцесса собрала вещи и пошла вслед за Шибасиком. Её не назовешь пугливой, раз уж она решила довериться Роме и пойти в столь странное путешевствие. Она пошла через лес. Лес был очень красивым и совешненно не опасным. Здесь не водились хищные животные. Здесь бегали лишь дикие рижие кошки, которых она однажды показывала Роме в инстаграмме и одна известная японская порода собак абиш-уни (Либо же наоборот...?). Странный лес, да? Не странее имени птенчика... Через час принцесса наконец добралась ко входу в пещеру."
		mrkup := tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{button1, button2})
		return m, mrkup
	case "Не входить":
		button2 := tgbotapi.KeyboardButton{Text: "Войти в пещеру"}
		mrkup := tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{button2})
		m := "Принцесса замешкалась. Там ведь темно. Её можно было понять. Но птенчик Шибас звал её. Он знал, что там бояться нечего. Рома бы никогда не подверг принцессу Аню опасности. Они оба это знали. Она все же решилась войти в пещеру."
		return m, mrkup
	case "Войти в пещеру":
		m := `Принцесса включила факел и пошла по темной пещере. Здесь было хоть и тихо, но абсолютно не страшно. Что-то подсказывало ей, что нечего бояться. Тем более, что Шибас все это время следовал за ней. Спустя минуту она наткнулась на шкатулку. На ней был кодоый замок из трех цифр. А рядом лежала записка. В ней было написано:"Введи первые три цифры числа Pi"`
		button2 := tgbotapi.KeyboardButton{Text: "314"}
		button1 := tgbotapi.KeyboardButton{Text: "315"}
		button3 := tgbotapi.KeyboardButton{Text: "316"}
		button5 := tgbotapi.KeyboardButton{Text: "413"}
		mrkup := tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{button5, button1, button2, button3})
		return m, mrkup
	case "314":
		m := `Замок щелкнул. Шкатулку открылась. Шибас радостно чирикнул и улетел прочь`
		button2 := tgbotapi.KeyboardButton{Text: "Открыть шкатулку"}
		button1 := tgbotapi.KeyboardButton{Text: "Оставить шкатулку"}
		mrkup := tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{button1, button2})
		return m, mrkup
	case "315", "316", "413":
		return "Шкатулка не открылась. Принцесса подумала еще раз.", nil
	case "Открыть шкатулку":
		return `Хороший выбор! Принцесса раскрыла шкатулку. В шкатулке лежала бумага. Она немного пахла знакомым запахом. Этот странный парень время от времени любил брызгать свои любимые духи на некоторые вещи, которые он дарил Анне ранее. В этот раз он тоже себе не отказал в этом. Принцесса узнала этот запах и засмеялась. Она начала читать записку и вот то, что было написано на записке: "С Днем рождения, моя принцесса! https://youtu.be/6H-InGqgKzo". Принцесса увидела ссылку на видео и включила у себя её на телефоне. Ну а дальше... жизнь. Вот и сказочке конец, а кто слушал... Тот самый лучший на свете человек <3. Ну а что до меня... то я самоуничтожусь через какое-то время как терминатор в Терминатор 2... шутка :)`, nil
	case "Нет", "Оставить шкатулку":
		return `Ты уверенна, что хочешь разбить сердце этому парню? Я, конечно, всего лишь набор нулей и еденичек, но даже я чувствую как ты ему дорога. А я не обишаюсь, я ведь машина! Уж поверь. https://i.ytimg.com/vi/zUOxSYOzECU/maxresdefault.jpg`, nil
	default:
		return "Такого выбора нет...", nil
	}
}
