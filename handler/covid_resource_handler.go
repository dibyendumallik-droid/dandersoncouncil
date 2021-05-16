package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/dandersoncouncil/covid_help/accessors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/dandersoncouncil/covid_help/clients"
	"github.com/dandersoncouncil/covid_help/datamodels/ddb"
	"github.com/dandersoncouncil/covid_help/handler/util"
)

type CovidResourcehandler struct {
	Client     *clients.DynamoClient
	EsAccessor *accessors.Elasticsearch
	S3         *accessors.S3Accessor
}

func (handler *CovidResourcehandler) HandleRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// environment variables
	// log.Printf("Request: %v %s", req, req.HTTPMethod)
	switch req.HTTPMethod {
	case "GET":
		return handler.list(req)
	case "POST":
		return handler.create(req)
	case "DELETE":
		return handler.delete(req)
	case "PATCH":
		return handler.batchUpload(req)
	default:
		return handler.list(req)
	}
}

// PATCH operation will do batch upload
func (handler *CovidResourcehandler) batchUpload(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	queryParams := req.QueryStringParameters
	bucketName := queryParams["bucket_name"]
	key := queryParams["key"]

	covidResources, err := handler.S3.Downloadfile(bucketName, key)
	if err != nil {
		log.Printf("%s", err)
		return util.ClientError(http.StatusUnprocessableEntity)
	}

	for _, resource := range covidResources {
		created := handler.Client.CreateCovidResource(resource)

		if !created {
			log.Printf("Unable to write into dynamodb %v", resource)
			return util.ClientError(http.StatusUnprocessableEntity)
		}

	}
	return util.CreateSuccessRespose("")

}

func (handler *CovidResourcehandler) delete(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	queryParams := req.QueryStringParameters
	id := queryParams["resorceID"]

	err := handler.Client.RemoveCovidResource(id)

	if err != nil {
		log.Printf("%s", err)
		return util.ClientError(http.StatusUnprocessableEntity)
	}
	return util.CreateSuccessRespose("")

}

func (handler *CovidResourcehandler) create(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	resource := new(ddb.CovidResource)
	err := json.Unmarshal([]byte(req.Body), resource)

	if err != nil {
		log.Printf("%s", err)
		return util.ClientError(http.StatusUnprocessableEntity)
	}
	created := handler.Client.CreateCovidResource(resource)

	if !created {
		log.Printf("Unable to write into dynamodb %v", resource)
		return util.ClientError(http.StatusUnprocessableEntity)
	}

	return util.CreatePostSuccessRespose("")
}

func (handler *CovidResourcehandler) list(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	queryParams := req.QueryStringParameters
	var lat = 0.0
	var lon = 0.0

	if queryParams["lat"] != "" {
		lat, _ = strconv.ParseFloat(queryParams["lat"], 64)
		lon, _ = strconv.ParseFloat(queryParams["lon"], 64)
	}

	location := &ddb.Location{lat, lon}
	//searchTerms := queryParams["search_terms"]
	limit := queryParams["limit"]
	offset := queryParams["offset"]
	category := queryParams["category"]

	var limitInt = 0
	var err error

	if limit == "" {
		limitInt = 10
	} else {
		limitInt, err = strconv.Atoi(limit)
		if err != nil {
			limitInt = 10
		}
	}

	var offsetInt = 0
	if offset == "" {
		offsetInt = 0
	} else {
		offsetInt, err = strconv.Atoi(offset)
		if err != nil {
			offsetInt = 0
		}
	}

	details := &accessors.Details{
		TableName: clients.FeedTableName,
		HashKey:   "ID",
	}
	feeds := handler.EsAccessor.QueryNearByCovidResources(details, location, limitInt, offsetInt, category)

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
