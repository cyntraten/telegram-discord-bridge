package telegram

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StartBot(token string) (*tgbotapi.BotAPI, error) {
	fmt.Println("Trying to connect to Telegram...")
	bot, err := tgbotapi.NewBotAPI(token)
	fmt.Println("Connected to Telegram!")
	if err != nil {
		return nil, err
	}

	return bot, nil
}

func StartListening(bot *tgbotapi.BotAPI) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.ChannelPost != nil {
			log.Printf("Новый пост в канале! Текст %v", update.ChannelPost.Text)

			// discord
		}
	}

}
