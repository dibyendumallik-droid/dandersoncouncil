package handler

import (
	"fmt"

	"log"

	"github.com/aws/aws-lambda-go/events"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/dandersoncouncil/covid_help/accessors"
	elastic "github.com/olivere/elastic/v7"
)

var awsSession = session.Must(session.NewSession(&aws.Config{}))
var dynamoSvc = dynamodb.New(awsSession)
var esclient = new(accessors.Elasticsearch)

type EsSyncLambda struct {
	SVC *elastic.Client
}

func (this *EsSyncLambda) Handler(e events.DynamoDBEvent) error {
	var item map[string]events.DynamoDBAttributeValue
	fmt.Println("Beginning ES Sync")
	for _, v := range e.Records {
		switch v.EventName {
		case "INSERT":
			fallthrough
		case "MODIFY":
			tableName := "Feed"
			item = v.Change.NewImage
			log.Printf("Processing item: %v", item)
			details := &accessors.Details{
				TableName: tableName,
				HashKey:   "ID",
			}
			esclient.Client = this.SVC
			resp, err := esclient.Update(details, item)
			if err != nil {
				log.Printf("Failure while calling ES: %v", err)
				return err
			}
			log.Printf("Response from ES: %v", resp.Result)
		default:
		}
	}
	return nil
}
