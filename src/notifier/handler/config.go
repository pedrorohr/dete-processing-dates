package handler

import (
	"log"
	"os"
	"strconv"
)

type config struct {
	deteBotApiToken string
	deteChatId      int64
}

func NewConfigFromEnv() *config {
	return &config{
		deteBotApiToken: os.Getenv("DETE_BOT_API_TOKEN"),
		deteChatId:      convertChatIdToInt64(os.Getenv("DETE_CHAT_ID")),
	}
}

func convertChatIdToInt64(chatId string) int64 {
	result, err := strconv.ParseInt(os.Getenv("DETE_CHAT_ID"), 10, 64)
	if err != nil {
		log.Fatalf("Got error converting Telegram chat ID to int64: %s", err)
	}

	return result
}
