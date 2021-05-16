package handler

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/dandersoncouncil/covid_help/clients"
	"github.com/dandersoncouncil/covid_help/datamodels/ddb"
	"github.com/dandersoncouncil/covid_help/handler/util"
	"github.com/google/uuid"
)

type CommentHandler struct {
	DDbClient *clients.DynamoClient
}

const (
	FeedIDQueryParam = "FeedID"
)

func (c *CommentHandler) HandleRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// environment variables
	// log.Printf("Request: %v %s", req, req.HTTPMethod)
	switch req.HTTPMethod {
	case "GET":
		return c.list(req)
	case "POST":
		return c.create(req)
	case "DELETE":
		return c.delete(req)
	default:
		return c.list(req)
	}
}

func (c *CommentHandler) list(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	queryParams := req.QueryStringParameters
	feedID := queryParams[FeedIDQueryParam]
	log.Printf("Query params: %v   all params: %v", feedID, queryParams)
	comments := c.DDbClient.ListComments(feedID)
	js, err := json.Marshal(comments)
	if err != nil {
		return util.ServerError(err)
	}
	return util.CreateHttpRespose(string(js))
}

func (c *CommentHandler) create(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	comments := new(ddb.Comments)
	err := json.Unmarshal([]byte(req.Body), comments)
	if err != nil {
		log.Printf("Error in CreateComment: %v", err)
	}
	commentID := uuid.New().String()
	user := util.GetUser(req)

	feedID := comments.FeedID
	userID := user.ID

	success := c.DDbClient.CreateComment(feedID, commentID, userID, comments.Data)

	if !success {
		return util.ServerError(errors.New("Cannot create comment"))
	}

	dbErr := c.DDbClient.UpdateCommentsCount(1, feedID)

	if dbErr != nil {
		return util.ServerError(dbErr)
	}

	return util.CreatePostSuccessRespose(commentID)
}

func (c *CommentHandler) delete(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	queryParams := req.QueryStringParameters
	feedID := queryParams["FeedID"]
	commentID := queryParams["CommentID"]

	err := c.DDbClient.RemoveComment(feedID, commentID)

	if err == nil {
		return util.ServerError(err)
	}
	return util.CreateSuccessRespose("OK")
}
