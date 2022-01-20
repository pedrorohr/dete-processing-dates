package utils

import (
	"log"
	"strings"
	"time"
)

const deteDateLayout = "02 January 2006"

func ExtractProcessingDate(htmlDate string) time.Time {
	dateWithoutNoBreakingSpaces := strings.ReplaceAll(htmlDate, "\xc2\xa0", " ")
	trimmedDate := strings.TrimSpace(dateWithoutNoBreakingSpaces)
	result, err := time.Parse(deteDateLayout, trimmedDate)
	if err != nil {
		log.Fatalf("Got error parsing processing date: %s", err)
	}

	return result
}
