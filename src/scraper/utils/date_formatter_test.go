package utils

import (
	"testing"
	"time"
)

type TestCase struct {
	name                   string
	htmlDate               string
	expectedProcessingDate time.Time
}

type TestCases []TestCase

func TestFormatDateForNotification(t *testing.T) {
	testCases := TestCases{
		{
			"RFC3339 example",
			"02 January 2006",
			getTime(2006, 1, 2),
		},
		{
			"May the Force be with you",
			"4 May 1977",
			getTime(1977, 5, 4),
		},
		{
			"Saint Patrick's Day",
			"17 March 2022",
			getTime(2022, 3, 17),
		},
		{
			"Date extracted from HTML - Day with 1 digit",
			"\xc2\xa01 April 2023\xc2\xa0",
			getTime(2023, 4, 1),
		},
		{
			"Date extracted from HTML - Day with 2 digits",
			"\xc2\xa017 March 2022\xc2\xa0",
			getTime(2022, 3, 17),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			formattedDate := ExtractProcessingDate(testCase.htmlDate)
			assertEquals(t, formattedDate, testCase.expectedProcessingDate)
		})
	}
}

func getTime(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func assertEquals(t testing.TB, got time.Time, want time.Time) {
	t.Helper()
	if got != want {
		t.Errorf("got %s, but want %s", got, want)
	}
}
