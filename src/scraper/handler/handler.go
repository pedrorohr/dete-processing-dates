package handler

type Handler interface {
	Run()
}

type lambdaHandler struct {
	deteUrl string
}

func (l lambdaHandler) Run() {

}

func NewLambdaHandler(deteUrl string) *lambdaHandler {
	return &lambdaHandler{
		deteUrl: deteUrl,
	}
}
