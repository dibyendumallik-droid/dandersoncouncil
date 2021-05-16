package handler

import (
	"encoding/json"
	"fmt"
	"strconv"

	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/dandersoncouncil/covid_help/accessors"
	"github.com/dandersoncouncil/covid_help/handler/util"
	"github.com/google/uuid"
)

type FeedMediaHandler struct {
	S3Accessor *accessors.S3Accessor
}

const (
	FeedMediaFolder = "feedMedia"
)

func (f *FeedMediaHandler) HandleRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	queryParams := req.QueryStringParameters
	feedID := queryParams["FeedID"]

	var urls map[string]string

	switch req.HTTPMethod {
	case "GET":
		mediaIDs := queryParams["mediaIDs"]
		log.Printf("Found media IDs %v", mediaIDs)
		urls = f.createGETURLs(feedID, mediaIDs)
	case "POST":
		mediaCount, _ := strconv.Atoi(queryParams["Count"])
		urls = f.createPOSTURLs(feedID, mediaCount)
	default:
		mediaIDs := queryParams["mediaIDs"]
		urls = f.createGETURLs(feedID, mediaIDs)
	}

	response, err := json.Marshal(urls)
	if err != nil {
		return util.ServerError(err)
	}
	return util.CreateHttpRespose(string(response))
}

func (f *FeedMediaHandler) createPOSTURLs(feedID string, count int) map[string]string {
	urls := make(map[string]string)
	for i := 0; i < count; i++ {
		mediaID := uuid.New().String()
		key := fmt.Sprintf("%v/%v", FeedMediaFolder, mediaID)
		url := f.S3Accessor.CreatePreSignedPutUrl(util.BUCKET_NAME, key)
		urls[mediaID] = url
	}
	return urls
}

func (f *FeedMediaHandler) createGETURLs(feedID string, mediaIDs string) map[string]string {
	urls := make(map[string]string)
	IDs := strings.Split(mediaIDs, ",")

	for _, v := range IDs {
		key := fmt.Sprintf("%v/%v", FeedMediaFolder, v)
		url := f.S3Accessor.CreatePreSignedGetUrl(util.BUCKET_NAME, key)
		urls[v] = url
		log.Printf("Pre-signed media url %v", url)
	}
	return urls
}
