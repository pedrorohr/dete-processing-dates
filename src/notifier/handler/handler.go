package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/pedrorohr/dete-processing-dates/src/notifier/utils"
)

const (
	parseMode           = "MarkdownV2"
	baseURL             = "https://api.telegram.org/bot"
	sendMessageEndpoint = "sendMessage"
	applicationType     = "application/json"
	baseNotification    = "A date has been updated!"
	briefcaseEmoji      = "💼"
	calendarEmoji       = "📅"
)

type Handler interface {
	Run(e events.DynamoDBEvent)
}

type lambdaHandler struct {
	deteBotApiToken string
	deteChatId      int64
}

type sendMessageReqBody struct {
	ChatId    int64  `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

func (l lambdaHandler) Run(e events.DynamoDBEvent) {
	for _, record := range e.Records {

		oldDateType := record.Change.OldImage["Type"].String()
		oldDate := record.Change.OldImage["Date"].String()
		log.Printf("Old - %s: %s\n", oldDateType, oldDate)

		newDateType := record.Change.NewImage["Type"].String()
		newDate := record.Change.NewImage["Date"].String()
		log.Printf("New - %s: %s\n", newDateType, newDate)

		text := getNotificationText(newDateType, newDate)
		l.sendTelegramMessage(text)
	}
}

func NewLambdaHandler(deteBotApiToken string, deteChatId int64) *lambdaHandler {
	return &lambdaHandler{
		deteBotApiToken: deteBotApiToken,
		deteChatId:      deteChatId,
	}
}

func (l lambdaHandler) sendTelegramMessage(text string) {
	reqBody := &sendMessageReqBody{
		ChatId:    l.deteChatId,
		Text:      text,
		ParseMode: parseMode,
	}
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.Fatalf("Got error marshalling sendMessage request body: %s", err)
	}

	uri := fmt.Sprintf("%s%s/%s", baseURL, l.deteBotApiToken, sendMessageEndpoint)
	res, err := http.Post(uri, applicationType, bytes.NewBuffer(reqBytes))
	if err != nil {
		log.Fatalf("Got error posting sendMessage request: %s", err)
	}
	if res.StatusCode != http.StatusOK {
		log.Fatalf("Got unexpected status after sending message to Telegram: %s", err)
	}
}

func getNotificationText(dateType, date string) string {
	return fmt.Sprintf("```\n%s\n\n%s %s\n\n%s %s\n```", baseNotification, briefcaseEmoji, dateType, calendarEmoji, utils.FormatDateForNotification(date))
}
