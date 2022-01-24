package handler

import (
	"github.com/pedrorohr/dete-processing-dates/src/scraper/dal"
)

func Create() Handler {
	config := NewConfigFromEnv()
	processingDateDal := dal.NewProcessingDate()

	return NewLambdaHandler(config.deteProcessingDatesUrl, processingDateDal)
}
