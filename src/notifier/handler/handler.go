package handler

import (
	"log"

	"github.com/aws/aws-lambda-go/events"

	"github.com/pedrorohr/dete-processing-dates/src/notifier/messenger"
)

type Handler interface {
	Run(e events.DynamoDBEvent)
}

type lambdaHandler struct {
	messenger messenger.Messenger
}

func NewLambdaHandler(messenger messenger.Messenger) *lambdaHandler {
	return &lambdaHandler{
		messenger: messenger,
	}
}

func (l lambdaHandler) Run(e events.DynamoDBEvent) {
	for _, record := range e.Records {

		oldDateType := record.Change.OldImage["Type"].String()
		oldDate := record.Change.OldImage["Date"].String()
		log.Printf("Old - %s: %s\n", oldDateType, oldDate)

		newDateType := record.Change.NewImage["Type"].String()
		newDate := record.Change.NewImage["Date"].String()
		log.Printf("New - %s: %s\n", newDateType, newDate)

		message := l.messenger.FormatMessage(newDateType, newDate)
		err := l.messenger.Dispatch(message)
		if err != nil {
			log.Fatalf("Got error dispatching message: %s", err)
		}
	}
}
