package telegram

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MessagePost struct {
	ChannelName  string
	Text         string
	FileURL      string
	FileName     string
	HasFile      bool
	MediaGroupID string
}

type MediaFile struct {
	URL      string
	FileName string
}

type AlbumPost struct {
	Text  string
	Files []MediaFile
}

func StartBot(token string) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return bot, nil
}

func StartListening(bot *tgbotapi.BotAPI, targetChannelID string) <-chan MessagePost {

	postsChan := make(chan MessagePost)

	go func() {
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60

		updates := bot.GetUpdatesChan(u)

		for update := range updates {

			if update.ChannelPost == nil {
				continue
			}

			currentChannelID := update.ChannelPost.Chat.ID
			if strconv.Itoa(int(currentChannelID)) != targetChannelID {
				continue
			}

			if update.ChannelPost != nil {

				post := MessagePost{
					MediaGroupID: update.ChannelPost.MediaGroupID,
					ChannelName:  update.ChannelPost.Chat.Title,
				}

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

type albumBuffer struct {
	text  string
	files []MediaFile
	timer *time.Timer
}

func StartAlbumCollection(inputChan <-chan MessagePost, outputChan chan<- AlbumPost) {
	buffer := make(map[string]*albumBuffer)
	var mu sync.Mutex

	for post := range inputChan {
		if post.MediaGroupID == "" {
			outputChan <- AlbumPost{
				Text: fmt.Sprintf("📢 **%s:**\n%s", post.ChannelName, post.Text),
				Files: []MediaFile{
					{URL: post.FileURL, FileName: post.FileName},
				},
			}
			continue
		}

		mu.Lock()
		buf, exists := buffer[post.MediaGroupID]

		if !exists {
			buf = &albumBuffer{
				text: post.Text,
				files: []MediaFile{
					{URL: post.FileURL, FileName: post.FileName},
				},
			}

			id := post.MediaGroupID
			buf.timer = time.AfterFunc(1*time.Second, func() {
				mu.Lock()
				completedAlbum := buffer[id]
				delete(buffer, id)
				mu.Unlock()

				outputChan <- AlbumPost{
					Text:  completedAlbum.text,
					Files: completedAlbum.files,
				}
			})

			buffer[id] = buf

		} else {
			buf.files = append(buf.files, MediaFile{
				URL:      post.FileURL,
				FileName: post.FileName,
			})

			if buf.text == "" && post.Text != "" {
				buf.text = post.Text
			}

			buf.timer.Reset(1 * time.Second)
		}

		mu.Unlock()
	}
}
