package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken    string
	DiscordToken     string
	DiscordChannelId string
}

func Load() (*Config, error) {

	_ = godotenv.Load()

	telegramToken := os.Getenv("TelegramToken")
	if telegramToken == "" {
		log.Fatal("TelegramToken is required, but it is empty")
	}

	discordToken := os.Getenv("DiscordToken")
	if discordToken == "" {
		log.Fatal("DiscordToken is required, but it is empty")
	}

	discordChannelId := os.Getenv("DiscordChannelId")
	if discordChannelId == "" {
		log.Fatal("DiscordChannelId is required, but it is empty")
	}

	return &Config{
		TelegramToken:    telegramToken,
		DiscordToken:     discordToken,
		DiscordChannelId: discordChannelId,
	}, nil

}
