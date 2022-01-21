package models

import "time"

type ProcessingDate struct {
	Type string
	Date time.Time
}

type ProcessingDates map[string]time.Time
