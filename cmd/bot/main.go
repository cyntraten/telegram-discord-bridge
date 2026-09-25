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

	rawsPostsChan := telegram.StartListening(tgbot, config.TargetTelegramChannelId)

	albumsChan := make(chan telegram.AlbumPost)
	go telegram.StartAlbumCollection(rawsPostsChan, albumsChan)

	for album := range albumsChan {
		var discordFiles []discord.MediaFileToSend

		for _, file := range album.Files {
			if file.URL == "" {
				continue
			}

			resp, err := http.Get(file.URL)
			if err != nil {
				log.Printf("Failed to download file from Telegram: %v", err)
				continue
			}

			discordFiles = append(discordFiles, discord.MediaFileToSend{Reader: resp.Body, FileName: file.FileName})
		}

		err = discord.SendAlbum(dsbot, config.DiscordChannelId, album.Text, discordFiles)
		if err != nil {
			log.Printf("Failed to send album in discord: %v", err)
		}

		for _, df := range discordFiles {
			if closer, ok := df.Reader.(io.ReadCloser); ok {
				closer.Close()
			}
		}

	}

}
