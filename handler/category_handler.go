package handler

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/dandersoncouncil/covid_help/handler/util"
)

type CategoryHandler struct {
}

func (handler *CategoryHandler) HandleRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	categories := []string{"Abulance", "Beds", "Doctors", "Food", "Helpline", "Medicine", "Oxygen", "Plasma", "Testing"}
	js, err := json.Marshal(categories)
	if err != nil {
		return util.ServerError(err)
	}
	return util.CreateHttpRespose(string(js))
}
