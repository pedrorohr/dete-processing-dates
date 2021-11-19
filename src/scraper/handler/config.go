package handler

import (
	"os"
)

type config struct {
	deteUrl string
}

// NewConfigFromEnv -
func NewConfigFromEnv() *config {

	return &config{
		deteUrl: os.Getenv("DETE_URL"),
	}
}
