package handler

import (
	"os"
)

type config struct {
	deteBotApiToken string
	deteChatId      string
}

func NewConfigFromEnv() *config {
	return &config{
		deteBotApiToken: os.Getenv("DETE_BOT_API_TOKEN"),
		deteChatId:      os.Getenv("DETE_CHAT_ID"),
	}
}
