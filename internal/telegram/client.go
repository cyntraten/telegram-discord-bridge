package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StartBot(token string) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return bot, nil
}

func StartListening(bot *tgbotapi.BotAPI) <-chan string {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	postsChan := make(chan string)

	go func() {
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60

		updates := bot.GetUpdatesChan(u)

		for update := range updates {
			if update.ChannelPost != nil {
				text := update.ChannelPost.Text
				if text != "" {
					log.Printf("New post in Telegram: %s", text)
					postsChan <- text
				}
			}
		}
	}()

	return postsChan

}
