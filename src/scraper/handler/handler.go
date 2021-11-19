package handler

import (
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

const (
	trustedPartnerDateSelector = "td:contains(\"Trusted Partner\") + td"
	standardDateSelector       = "td:contains(\"Standard\") + td"
	deteDateLayout             = "02 January 2006"
)

type Handler interface {
	Run()
}

type lambdaHandler struct {
	deteProcessingDatesUrl string
}

func (l lambdaHandler) Run() {
	c := colly.NewCollector()

	c.OnHTML(trustedPartnerDateSelector, evaluateTrustedPartnerProcessingDate)

	c.OnHTML(standardDateSelector, evaluateStandardProcessingDate)

	c.Visit(l.deteProcessingDatesUrl)
}

func NewLambdaHandler(deteProcessingDatesUrl string) *lambdaHandler {
	return &lambdaHandler{
		deteProcessingDatesUrl: deteProcessingDatesUrl,
	}
}

func evaluateTrustedPartnerProcessingDate(e *colly.HTMLElement) {
	extracted := extractProcessingDate(e.Text)

	fmt.Printf("\nPartner: %s", extracted)
}

func evaluateStandardProcessingDate(e *colly.HTMLElement) {
	extracted := extractProcessingDate(e.Text)

	fmt.Printf("\nStandard: %s", extracted)
}

func extractProcessingDate(htmlDate string) time.Time {
	dateWithoutNoBreakingSpaces := strings.ReplaceAll(htmlDate, "\xc2\xa0", " ")
	trimmedDate := strings.TrimSpace(dateWithoutNoBreakingSpaces)
	result, err := time.Parse(deteDateLayout, trimmedDate)
	if err != nil {
		panic(err)
	}

	return result
}
