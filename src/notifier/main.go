package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pedrorohr/dete-processing-dates/src/notifier/handler"
)

func main() {
	handler := handler.Create()
	lambda.Start(handler.Run)
}
