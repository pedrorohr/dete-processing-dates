package handler

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
)

type Handler interface {
	Run(e events.DynamoDBEvent)
}

type lambdaHandler struct {
	deteBotApiToken string
	deteChatId      string
}

func (l lambdaHandler) Run(e events.DynamoDBEvent) {
	for _, record := range e.Records {

		oldDateType := record.Change.OldImage["Type"].String()
		oldDate := record.Change.OldImage["Date"].String()
		log.Printf("Old - Type: %s Date: %s\n", oldDateType, oldDate)

		newDateType := record.Change.NewImage["Type"].String()
		newDate := record.Change.NewImage["Date"].String()
		log.Printf("New - Type: %s Date: %s\n", newDateType, newDate)
	}
}

func NewLambdaHandler(deteBotApiToken string, deteChatId string) *lambdaHandler {
	return &lambdaHandler{
		deteBotApiToken: deteBotApiToken,
		deteChatId:      deteChatId,
	}
}
