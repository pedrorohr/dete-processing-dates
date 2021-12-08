package handler

import (
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gocolly/colly"
)

const (
	trustedPartnerDateSelector = "td:contains(\"Trusted Partner\") + td"
	standardDateSelector       = "td:contains(\"Standard\") + td"
	deteDateLayout             = "02 January 2006"
	tableName                  = "DeteProcessingDates"
	trusted                    = "Trusted Partner"
	standard                   = "Standard"
)

type Handler interface {
	Run()
}

type lambdaHandler struct {
	deteProcessingDatesUrl string
	db                     *dynamodb.DynamoDB
}

type deteProcessingDate struct {
	Type string
	Date time.Time
}

type deteProcessingDates map[string]time.Time

func (l lambdaHandler) Run() {
	dates := loadDeteProcessingDates(l.db)

	c := colly.NewCollector()

	c.OnHTML(trustedPartnerDateSelector, func(e *colly.HTMLElement) {
		l.evaluateProcessingDate(e, trusted, dates)
	})

	c.OnHTML(standardDateSelector, func(e *colly.HTMLElement) {
		l.evaluateProcessingDate(e, standard, dates)
	})

	c.Visit(l.deteProcessingDatesUrl)
}

func NewLambdaHandler(deteProcessingDatesUrl string, db *dynamodb.DynamoDB) *lambdaHandler {
	return &lambdaHandler{
		deteProcessingDatesUrl: deteProcessingDatesUrl,
		db:                     db,
	}
}

func (l lambdaHandler) evaluateProcessingDate(e *colly.HTMLElement, dateType string, dates deteProcessingDates) {
	extractedDate := extractProcessingDate(e.Text)
	log.Printf("Extracted date - %s: %s\n", dateType, extractedDate)

	actualDate, ok := dates[dateType]
	if !ok || extractedDate.After(actualDate) {
		newActualDate := deteProcessingDate{
			Type: dateType,
			Date: extractedDate,
		}
		l.saveNewProcessingDate(newActualDate)
	}
}

func (l lambdaHandler) saveNewProcessingDate(date deteProcessingDate) {
	avDate, err := dynamodbattribute.MarshalMap(date)
	if err != nil {
		log.Fatalf("Got error marshalling new date: %s", err)
	}

	input := &dynamodb.PutItemInput{
		Item:      avDate,
		TableName: aws.String(tableName),
	}

	_, err = l.db.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
	}
}

func loadDeteProcessingDates(db *dynamodb.DynamoDB) deteProcessingDates {
	params := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	result, err := db.Scan(params)
	if err != nil {
		log.Fatalf("Scan API call failed: %s", err)
	}

	deteProcessingDates := make(map[string]time.Time)
	for _, i := range result.Items {
		item := deteProcessingDate{}

		err = dynamodbattribute.UnmarshalMap(i, &item)
		if err != nil {
			log.Fatalf("Got error unmarshalling: %s", err)
		}

		log.Printf("Actual date - %s: %s\n", item.Type, item.Date)
		deteProcessingDates[item.Type] = item.Date
	}

	return deteProcessingDates
}

func extractProcessingDate(htmlDate string) time.Time {
	dateWithoutNoBreakingSpaces := strings.ReplaceAll(htmlDate, "\xc2\xa0", " ")
	trimmedDate := strings.TrimSpace(dateWithoutNoBreakingSpaces)
	result, err := time.Parse(deteDateLayout, trimmedDate)
	if err != nil {
		log.Fatalf("Got error parsing processing date: %s", err)
	}

	return result
}
