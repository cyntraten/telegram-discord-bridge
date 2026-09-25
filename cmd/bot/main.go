package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

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
	for post := range postsChan {

		var fileReader io.Reader = nil
		var resp *http.Response
		var err error

		if post.FileURL != "" && post.HasFile != false {
			resp, err = http.Get(post.FileURL)
			if err == nil {
				fileReader = resp.Body
			} else {
				log.Printf("Failed to download image from Telegram post: %v\n", err)
			}
		}

		err = discord.SendMessage(dsbot, config.DiscordChannelId, post.Text, post.FileName, fileReader, post.HasFile)
		if err != nil {
			log.Printf("Failed to send message in discord %v", err)
		}

		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}

}
