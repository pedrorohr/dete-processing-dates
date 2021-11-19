package handler

func Create() Handler {
	config := NewConfigFromEnv()

	return NewLambdaHandler(config.deteProcessingDatesUrl)
}
