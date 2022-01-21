package dal

import (
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/pedrorohr/dete-processing-dates/src/scraper/models"
)

const tableName = "DeteProcessingDates"

type ProcessingDateDAL interface {
	LoadAll() models.ProcessingDates
	Save(processingDate models.ProcessingDate)
}

type ProcessingDate struct {
	db *dynamodb.DynamoDB
}

func NewProcessingDate() *ProcessingDate {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	db := dynamodb.New(sess)

	return &ProcessingDate{
		db: db,
	}
}

func (pd *ProcessingDate) LoadAll() models.ProcessingDates {
	params := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	result, err := pd.db.Scan(params)
	if err != nil {
		log.Fatalf("Scan API call failed: %s", err)
	}

	deteProcessingDates := make(map[string]time.Time)
	for _, i := range result.Items {
		item := models.ProcessingDate{}

		err = dynamodbattribute.UnmarshalMap(i, &item)
		if err != nil {
			log.Fatalf("Got error unmarshalling: %s", err)
		}

		log.Printf("Actual date - %s: %s\n", item.Type, item.Date)
		deteProcessingDates[item.Type] = item.Date
	}

	return deteProcessingDates
}

func (pd *ProcessingDate) Save(processingDate models.ProcessingDate) {
	avDate, err := dynamodbattribute.MarshalMap(processingDate)
	if err != nil {
		log.Fatalf("Got error marshalling new date: %s", err)
	}

	input := &dynamodb.PutItemInput{
		Item:      avDate,
		TableName: aws.String(tableName),
	}

	_, err = pd.db.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
	}
}
