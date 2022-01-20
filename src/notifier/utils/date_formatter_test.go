package utils

import "testing"

type TestCase struct {
	name                  string
	fullDateTime          string
	expectedFormattedDate string
}

type TestCases []TestCase

func TestFormatDateForNotification(t *testing.T) {
	testCases := TestCases{
		{
			"RFC3339 example",
			"2006-01-02T15:04:05Z",
			"02 January 2006",
		},
		{
			"Saint Patrick's Day",
			"2022-03-17T23:53:55Z",
			"17 March 2022",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			formattedDate := FormatDateForNotification(testCase.fullDateTime)
			assertEquals(t, formattedDate, testCase.expectedFormattedDate)
		})
	}
}

func assertEquals(t testing.TB, got string, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %s, but want %s", got, want)
	}
}
