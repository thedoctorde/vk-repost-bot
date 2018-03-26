package tg

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

type Bot struct {
	NewMessage string
	bot        *tgbotapi.BotAPI
}

func NewBot(token string) *Bot {
	b := &Bot{}
	var err error
	b.bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	b.bot.Debug = true
	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	return b
}

func (b *Bot) SendMessage(tgChannelName, message string) error {
	msg := tgbotapi.NewMessageToChannel(tgChannelName, message)
	_, err := b.bot.Send(msg)
	return err
}
