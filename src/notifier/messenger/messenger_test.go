package messenger

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

const (
	botApiToken = "1234"
	chatId      = int64(4321)
	message     = "Hello World!"
)

type FormatMessageTestCase struct {
	name            string
	dateType        string
	date            string
	expectedMessage string
}

type FormatMessageTestCases []FormatMessageTestCase

type DispatchTestCase struct {
	name          string
	httpStatus    int
	expectedError error
}

type DispatchTestCases []DispatchTestCase

func TestFormatMessage(t *testing.T) {
	testCases := FormatMessageTestCases{
		{
			"New Year's Eve",
			"Trusted Partner",
			"2021-12-31T23:59:59Z",
			"```\nA date has been updated!\n\n💼 Trusted Partner\n\n📅 31 December 2021\n```",
		},
		{
			"Saint Patrick's Day",
			"Standard",
			"2022-03-17T23:53:55Z",
			"```\nA date has been updated!\n\n💼 Standard\n\n📅 17 March 2022\n```",
		},
	}
	messenger, _ := NewTelegram(botApiToken, chatId)

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			formattedDate := messenger.FormatMessage(testCase.dateType, testCase.date)
			assertEquals(t, formattedDate, testCase.expectedMessage)
		})
	}
}

func TestDispatch(t *testing.T) {
	testCases := DispatchTestCases{
		{
			"Happy path",
			http.StatusOK,
			nil,
		},
		{
			"Bad request",
			http.StatusBadRequest,
			errors.New("got unexpected status after sending message to Telegram: 400 Bad Request"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			server := makeServer(testCase.httpStatus)
			defer server.Close()

			url, _ := url.Parse(server.URL)
			messenger, _ := NewTelegram(botApiToken, chatId, WithBaseURL(*url))

			err := messenger.Dispatch(message)
			assertError(t, err, testCase.expectedError)
		})
	}

}

func assertEquals(t testing.TB, got string, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %s, but want %s", got, want)
	}
}

func makeServer(httpStatusCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(httpStatusCode)
		rw.Write(nil)
	}))
}

func assertError(t testing.TB, got error, want error) {
	t.Helper()
	if got != nil && want == nil {
		t.Fatal("got an error but didn't want one")
	}
	if got == nil && want != nil {
		t.Fatal("didn't get an error but wanted one")
	}

	if got != nil && want != nil && got.Error() != want.Error() {
		t.Errorf("got %s, want %s", got, want)
	}
}
