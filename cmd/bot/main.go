package main

import (
	"fmt"
	"log"

	"github.com/cyntraten/telegram-discord-bridge/internal/config"
	"github.com/cyntraten/telegram-discord-bridge/internal/telegram"
)

func main() {
	config, err := config.Load()

	if err != nil {
		log.Fatal(err)
	}

	bot, err := telegram.StartBot(config.TelegramToken)

	if err != nil {
		log.Fatalf("Failed to start telegram bot: %v", err)
	}

	fmt.Printf("Bot started success\n")
	//bot.Debug = true

	telegram.StartListening(bot)

}
