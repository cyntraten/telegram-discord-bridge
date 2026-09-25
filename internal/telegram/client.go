package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MessagePost struct {
	Text         string
	FileURL      string
	FileName     string
	HasFile      bool
	MediaGroupID string
}

func StartBot(token string) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return bot, nil
}

func StartListening(bot *tgbotapi.BotAPI) <-chan MessagePost {
	postsChan := make(chan MessagePost)

	go func() {
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60

		updates := bot.GetUpdatesChan(u)

		for update := range updates {
			if update.ChannelPost != nil {

				post := MessagePost{}

				//text post
				if update.ChannelPost.Text != "" {
					post.Text = update.ChannelPost.Text
					log.Printf("New message in telegram channel: %v\n", post.Text)
				}

				var fileID string

				//post with image
				if len(update.ChannelPost.Photo) > 0 {
					photos := update.ChannelPost.Photo
					bestPhoto := photos[len(photos)-1]
					fileID = bestPhoto.FileID
					post.FileName = "image.png"
				}

				//post with video
				if update.ChannelPost.Video != nil {
					fileID = update.ChannelPost.Video.FileID
					post.FileName = "video.mp4"

				}

				//post with gif
				if update.ChannelPost.Animation != nil {
					fileID = update.ChannelPost.Animation.FileID
					post.FileName = "animation.mp4"
				}

				//post witch sticker
				if update.ChannelPost.Sticker != nil {
					fileID = update.ChannelPost.Sticker.FileID
					if update.ChannelPost.Sticker.IsAnimated {
						post.FileName = "sticker.webm"
					} else {
						post.FileName = "sticker.webp"
					}
				}

				//has file
				if fileID != "" {
					fileURL, err := bot.GetFileDirectURL(fileID)
					if err != nil {
						log.Printf("Failed to load file url: %v\n", err)
					} else {
						post.FileURL = fileURL
						post.HasFile = true
					}
				}

				if update.ChannelPost.Caption != "" {
					post.Text = update.ChannelPost.Caption
				}

				if post.HasFile {
					log.Printf("New message with attachment (%s) in telegram channel \n", post.Text)
				}

				if post.Text != "" || post.HasFile {
					postsChan <- post
				}
			}

		}
	}()

	return postsChan

}
