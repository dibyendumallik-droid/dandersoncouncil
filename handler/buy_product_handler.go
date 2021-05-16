package handler

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/dandersoncouncil/covid_help/clients"
	"github.com/dandersoncouncil/covid_help/datamodels/ddb"
	"github.com/dandersoncouncil/covid_help/handler/util"
	"github.com/google/uuid"
)

type BuyProductHandlerHandler struct {
	ddb *clients.DynamoClient
}

func (b *BuyProductHandlerHandler) HandleRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// environment variables
	// log.Printf("Request: %v %s", req, req.HTTPMethod)
	switch req.HTTPMethod {
	case "GET":
	case "POST":
		return b.create(req)
	case "DELETE":
	default:
		return b.create(req)
	}

	return b.create(req)
}

func (b *BuyProductHandlerHandler) create(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	orderId := uuid.NewString()

	order := new(ddb.Order)
	err := json.Unmarshal([]byte(req.Body), order)

	if err != nil {
		log.Printf("Error in CreateComment: %v", err)
		return util.ServerError(err)
	}
	user := util.GetUser(req)

	order.OrderId = orderId
	order.CustomerId = user.ID
	return util.ServerError(err)
}
