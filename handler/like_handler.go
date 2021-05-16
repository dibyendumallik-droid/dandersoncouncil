package handler

import (
	"encoding/json"

	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/dandersoncouncil/covid_help/clients"
	"github.com/dandersoncouncil/covid_help/datamodels/ddb"
	"github.com/dandersoncouncil/covid_help/handler/util"
)

type LikeHandler struct {
	DDbClient *clients.DynamoClient
}

func (l *LikeHandler) HandleRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// environment variables
	// log.Printf("Request: %v %s", req, req.HTTPMethod)
	switch req.HTTPMethod {
	case "POST":
		return l.create(req)
	case "DELETE":
		return l.delete(req)
	default:
		return l.create(req)
	}
}

func (l *LikeHandler) create(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	like := getLike(req)
	err := l.DDbClient.CreateLike(like)

	if err != nil {
		return util.ServerError(err)
	}

	err = l.DDbClient.UpdateLikeCount(1, like.FeedID)

	if err != nil {
		return util.ServerError(err)
	}
	return util.CreatePostSuccessRespose("OK")
}

func (l *LikeHandler) delete(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	like := getLike(req)
	err := l.DDbClient.RemoveLike(like)

	if err != nil {
		return util.ServerError(err)
	}

	err = l.DDbClient.UpdateLikeCount(-1, like.FeedID)

	if err != nil {
		return util.ServerError(err)
	}
	return util.CreateSuccessRespose("OK")
}

func getLike(req events.APIGatewayProxyRequest) *ddb.Like {
	like := new(ddb.Like)
	err := json.Unmarshal([]byte(req.Body), like)
	if err != nil {
		log.Printf("Error in parsing payload: %v", err)
	}
	return like
}
