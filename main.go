package main

import (
	"fmt"
	"log"
	"os"

	elastic "github.com/olivere/elastic/v7"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/dandersoncouncil/covid_help/accessors"
	"github.com/dandersoncouncil/covid_help/clients"
	"github.com/dandersoncouncil/covid_help/handler"
	"github.com/dandersoncouncil/covid_help/router"
)

const (
	STATIC_DIR = "/static/"
)

func ddbStreamHandler() {
	svc, err := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
		elastic.SetURL(fmt.Sprintf("https://%s", os.Getenv("ELASTICSEARCH_URL"))),
	)
	if err != nil {
		log.Printf("Error while creating ES client: %v", err)
		return
	}
	esSync := handler.EsSyncLambda{svc}
	lambda.Start(esSync.Handler)
}

func createEsClient() *accessors.Elasticsearch {
	elasticClient, _ := elastic.NewClient(
		elastic.SetHealthcheck(false),
		elastic.SetSniff(false),
		elastic.SetURL(fmt.Sprintf("https://%s", os.Getenv("ELASTICSEARCH_URL"))),
	)
	var esclient = new(accessors.Elasticsearch)
	esclient.Client = elasticClient
	return esclient
}
func apiGatewayHandler() {
	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Initialze AWS resources
	svc := dynamodb.New(sess)
	s3 := s3.New(sess)
	esclient := createEsClient()

	//Initialize accessors
	dynamoClient := clients.DynamoClient{svc}
	s3Accessor := accessors.S3Accessor{s3}

	// Initialize the handlers
	feedHandler := handler.FeedHandler{Client: &dynamoClient, EsAccessor: esclient}
	picHandler := handler.PicHandler{&s3Accessor}
	userHandler := handler.UserHandler{&s3Accessor, &dynamoClient}
	commentHandler := handler.CommentHandler{&dynamoClient}
	feedStatHandler := handler.FeedStatHandler{&dynamoClient}
	feedMediaHandler := handler.FeedMediaHandler{&s3Accessor}
	covidResourceHandler := handler.CovidResourcehandler{Client: &dynamoClient, EsAccessor: esclient, S3: &s3Accessor}
	categoryHandler := handler.CategoryHandler{}

	globalRouter := router.GlobalRouter{&feedHandler, &picHandler,
		&userHandler, &commentHandler, &feedStatHandler, &feedMediaHandler, &covidResourceHandler, &categoryHandler}

	lambda.Start(globalRouter.HandleRequest)
}

func main() {
	apiGatewayHandler()
}
