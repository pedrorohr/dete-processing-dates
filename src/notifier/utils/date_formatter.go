package utils

import (
	"log"
	"time"
)

const dateFormat = "02 January 2006"

func FormatDateForNotification(fullDateTime string) string {
	date, err := time.Parse(time.RFC3339, fullDateTime)

	if err != nil {
		log.Fatalf("Got error parsing processing date: %s", err)
	}

	return date.Format(dateFormat)
}
