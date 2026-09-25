package main

import (
	"fmt"
	"log"

	"github.com/cyntraten/telegram-discord-bridge/internal/config"
	"github.com/cyntraten/telegram-discord-bridge/internal/discord"
	"github.com/cyntraten/telegram-discord-bridge/internal/telegram"
)

func main() {
	config, err := config.Load()

	if err != nil {
		log.Fatal(err)
	}

	//init tgbot
	tgbot, err := telegram.StartBot(config.TelegramToken)

	if err != nil {
		log.Fatalf("Failed to start telegram bot: %v", err)
	}

	fmt.Printf("Bot started success\n")

	//init discordbot
	dsbot, err := discord.StartBot(config.DiscordToken)
	if err != nil {
		log.Fatalf("Failed to start discord bot: %v", err)
	}
	defer dsbot.Close()

	// channel with posts from tg
	postsChan := telegram.StartListening(tgbot)

	// send posts from tg in discord
	for text := range postsChan {
		discord.SendMessage(dsbot, config.DiscordChannelId, text)
	}

}
