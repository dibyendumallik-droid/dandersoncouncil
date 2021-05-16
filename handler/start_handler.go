package handler

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/dandersoncouncil/covid_help/clients"
	"github.com/dandersoncouncil/covid_help/handler/util"
)

type FeedStatHandler struct {
	DDbClient *clients.DynamoClient
}

func (l *FeedStatHandler) HandleRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	queryParams := req.QueryStringParameters

	feedID := queryParams["FeedID"]

	response, err := l.DDbClient.GetFeedStat(feedID)

	if err != nil {
		return util.ServerError(err)
	}

	js, err := json.Marshal(response)
	if err != nil {
		return util.ServerError(err)
	}

	// Return a response with a 200 OK status and the JSON book record
	// as the body.
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(js),
	}, nil
}
