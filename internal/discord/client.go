package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func StartBot(token string) (*discordgo.Session, error) {
	session, err := discordgo.New("Bot " + token)

	if err != nil {
		return nil, err
	}

	err = session.Open()

	if err != nil {
		return nil, err
	}

	return session, nil
}

func SendMessage(session *discordgo.Session, channelId string, textMessage string) {
	_, err := session.ChannelMessageSend(channelId, textMessage)
	if err != nil {
		fmt.Printf("Send message in channel: %v, error: %v\n", channelId, err)
	}
}
