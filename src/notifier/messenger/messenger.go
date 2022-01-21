package messenger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/pedrorohr/dete-processing-dates/src/notifier/utils"
)

const (
	parseMode           = "MarkdownV2"
	defaultBaseURL      = "https://api.telegram.org"
	sendMessageEndpoint = "sendMessage"
	applicationType     = "application/json"
	baseNotification    = "A date has been updated!"
	briefcaseEmoji      = "💼"
	calendarEmoji       = "📅"
)

type Messenger interface {
	FormatMessage(dateType, date string) string
	Dispatch(message string) error
}

type Telegram struct {
	botApiToken string
	chatId      int64
	baseURL     url.URL
}

type sendMessageReqBody struct {
	ChatId    int64  `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

type Option func(*Telegram)

func NewTelegram(botApiToken string, chatId int64, opts ...Option) (*Telegram, error) {
	baseUrl, err := url.Parse(defaultBaseURL)
	if err != nil {
		return nil, err
	}

	telegram := &Telegram{
		botApiToken: botApiToken,
		chatId:      chatId,
		baseURL:     *baseUrl,
	}

	for _, opt := range opts {
		opt(telegram)
	}

	return telegram, nil
}

func WithBaseURL(baseURL url.URL) Option {
	return func(t *Telegram) {
		t.baseURL = baseURL
	}
}

func (t *Telegram) FormatMessage(dateType, date string) string {
	return fmt.Sprintf("```\n%s\n\n%s %s\n\n%s %s\n```", baseNotification, briefcaseEmoji, dateType, calendarEmoji, utils.FormatDateForNotification(date))
}

func (t *Telegram) Dispatch(message string) error {
	reqBody := &sendMessageReqBody{
		ChatId:    t.chatId,
		Text:      message,
		ParseMode: parseMode,
	}
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("got error marshalling sendMessage request body: %s", err)
	}

	uri := fmt.Sprintf("%s/bot%s/%s", t.baseURL.String(), t.botApiToken, sendMessageEndpoint)
	res, err := http.Post(uri, applicationType, bytes.NewBuffer(reqBytes))
	if err != nil {
		return fmt.Errorf("got error posting sendMessage request: %s", err)
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("got unexpected status after sending message to Telegram: %s", res.Status)
	}

	return nil
}
