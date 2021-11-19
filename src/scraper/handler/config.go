package handler

import (
	"os"
)

type config struct {
	deteProcessingDatesUrl string
}

// NewConfigFromEnv -
func NewConfigFromEnv() *config {

	return &config{
		deteProcessingDatesUrl: os.Getenv("DETE_PROCESSING_DATES_URL"),
	}
}
