package handler

import "github.com/pedrorohr/dete-processing-dates/src/notifier/messenger"

func Create() Handler {
	config := NewConfigFromEnv()

	messenger, _ := messenger.NewTelegram(config.deteBotApiToken, config.deteChatId)
	return NewLambdaHandler(messenger)
}
