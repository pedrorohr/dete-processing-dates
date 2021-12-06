package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

const (
	baseURL             = "https://api.telegram.org/bot"
	sendMessageEndpoint = "sendMessage"
	applicationType     = "application/json"
	baseNotification    = "A date has been updated!"
	briefcaseEmoji      = "💼"
	calendarEmoji       = "📅"
	dateFormat          = "2 January 2006"
)

type Handler interface {
	Run(e events.DynamoDBEvent)
}

type lambdaHandler struct {
	deteBotApiToken string
	deteChatId      int64
}

type sendMessageReqBody struct {
	ChatId int64  `json:"chat_id"`
	Text   string `json:"text"`
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
		ChatId: l.deteChatId,
		Text:   text,
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
	return fmt.Sprintf("```\n%s\n\n%s %s\n\n%s %s\n```", baseNotification, briefcaseEmoji, dateType, calendarEmoji, formatDate(date))
}

func formatDate(fullDateTime string) string {
	date, err := time.Parse(time.RFC3339, fullDateTime)
	if err != nil {
		log.Fatalf("Got error parsing processing date: %s", err)
	}

	return date.Format(dateFormat)
}
