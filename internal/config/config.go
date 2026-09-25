package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken string
	DiscordToken  string
}

func Load() (*Config, error) {

	_ = godotenv.Load()

	telegramToken := os.Getenv("TelegramToken")
	discordToken := os.Getenv("DiscordToken")

	if telegramToken == "" {
		log.Fatal("TelegramToken is required, but it is empty")
	}

	if discordToken == "" {
		log.Fatal("DiscordToken is required, but it is empty")
	}

	return &Config{
		TelegramToken: telegramToken,
		DiscordToken:  discordToken,
	}, nil

}
