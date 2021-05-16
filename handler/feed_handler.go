package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/dandersoncouncil/covid_help/accessors"
	"github.com/dandersoncouncil/covid_help/clients"
	"github.com/dandersoncouncil/covid_help/datamodels/ddb"
	"github.com/dandersoncouncil/covid_help/handler/util"
	"github.com/google/uuid"
)

type FeedHandler struct {
	Client     *clients.DynamoClient
	EsAccessor *accessors.Elasticsearch
}

func (this *FeedHandler) HandleRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// environment variables
	// log.Printf("Request: %v %s", req, req.HTTPMethod)
	switch req.HTTPMethod {
	case "GET":
		return this.list(req)
	case "POST":
		return this.create(req)
	case "DELETE":
		return this.delete(req)
	case "PATCH":
		return this.patch(req)
	default:
		return this.list(req)
	}
}

func (this *FeedHandler) list(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	details := &accessors.Details{
		TableName: clients.FeedTableName,
		HashKey:   "ID",
	}

	queryParams := req.QueryStringParameters
	var lat = 0.0
	var lon = 0.0

	if queryParams["lat"] != "" {
		lat, _ = strconv.ParseFloat(queryParams["lat"], 64)
		lon, _ = strconv.ParseFloat(queryParams["lon"], 64)
	}

	location := &ddb.Location{lat, lon} //TODO Get this from args

	feeds := this.EsAccessor.QueryNearByFeedElement(details, location)
	js, err := json.Marshal(feeds)
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

func (this *FeedHandler) create(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	feed := new(ddb.Feed)
	err := json.Unmarshal([]byte(req.Body), feed)
	id := uuid.New().String()

	now := time.Now()
	secs := now.Unix()

	log.Printf("Created id: %s", id)
	if err != nil {
		log.Printf("%s", err)
		return util.ClientError(http.StatusUnprocessableEntity)
	}

	context := req.RequestContext
	claims := context.Authorizer["claims"].(map[string]interface{})
	name := fmt.Sprintf("%v", claims["name"])

	feed.ID = id
	feed.CreatedBy = name
	feed.CreationTime = secs
	this.Client.CreateFeed(feed)
	return util.CreatePostSuccessRespose(id)
}

func (f *FeedHandler) patch(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	feed := new(ddb.Feed)
	err := json.Unmarshal([]byte(req.Body), feed)

	if err != nil {
		return util.ServerError(err)
	}

	log.Printf("PATCHing feed %v with list %v", feed.ID, feed.MediaList)
	err = f.Client.PathFeed(feed.ID, feed.MediaList)

	if err != nil {
		return util.ServerError(err)
	}

	return util.CreatePostSuccessRespose(feed.ID)

}

func (f *FeedHandler) delete(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	queryParams := req.QueryStringParameters
	feedID := queryParams["FeedID"]
	err := f.Client.RemoveFeed(feedID)

	if err == nil {
		return util.ServerError(err)
	}
	return util.CreateSuccessRespose("OK")
}
