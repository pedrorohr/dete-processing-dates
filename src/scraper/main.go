package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pedrorohr/dete-processing-dates/src/scraper/handler"
)

func main() {
	handler := handler.Create()
	lambda.Start(handler.Run)
}
