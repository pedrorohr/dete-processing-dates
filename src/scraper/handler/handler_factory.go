package handler

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func Create() Handler {
	config := NewConfigFromEnv()

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	db := dynamodb.New(sess)

	return NewLambdaHandler(config.deteProcessingDatesUrl, db)
}
