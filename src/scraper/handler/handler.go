package handler

import (
	"log"

	"github.com/gocolly/colly"

	"github.com/pedrorohr/dete-processing-dates/src/scraper/dal"
	"github.com/pedrorohr/dete-processing-dates/src/scraper/models"
	"github.com/pedrorohr/dete-processing-dates/src/scraper/utils"
)

const (
	trustedPartnerDateSelector = "td:contains(\"Trusted Partner\") + td"
	standardDateSelector       = "td:contains(\"Standard\") + td"
	trusted                    = "Trusted Partner"
	standard                   = "Standard"
)

func NewLambdaHandler(deteProcessingDatesUrl string, processingDateDal dal.ProcessingDateDAL) *lambdaHandler {
	return &lambdaHandler{
		deteProcessingDatesUrl: deteProcessingDatesUrl,
		processingDateDal:      processingDateDal,
	}
}

type Handler interface {
	Run()
}

type lambdaHandler struct {
	deteProcessingDatesUrl string
	processingDateDal      dal.ProcessingDateDAL
}

func (l lambdaHandler) Run() {
	dates := l.processingDateDal.LoadAll()

	c := colly.NewCollector()

	c.OnHTML(trustedPartnerDateSelector, func(e *colly.HTMLElement) {
		l.evaluateProcessingDate(e, trusted, dates)
	})

	c.OnHTML(standardDateSelector, func(e *colly.HTMLElement) {
		l.evaluateProcessingDate(e, standard, dates)
	})

	c.Visit(l.deteProcessingDatesUrl)
}

func (l lambdaHandler) evaluateProcessingDate(e *colly.HTMLElement, dateType string, dates models.ProcessingDates) {
	extractedDate := utils.ExtractProcessingDate(e.Text)
	log.Printf("Extracted date - %s: %s\n", dateType, extractedDate)

	actualDate, ok := dates[dateType]
	if !ok || extractedDate.After(actualDate) {
		newActualDate := models.ProcessingDate{
			Type: dateType,
			Date: extractedDate,
		}
		l.processingDateDal.Save(newActualDate)
	}
}
