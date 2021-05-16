package builders

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/dandersoncouncil/covid_help/clients"
)

var dynamoInitialized uint32
var s3ClientInitialized uint32

var dynamo *clients.DynamoClient
var s3client *clients.S3Client

// Returns a singleton instance of dynamo client
func GetDynamoInstance() *clients.DynamoClient {

	if atomic.LoadUInt32(&dynamoInitialized) == 1 {
		return dynamo
	}

	mu.Lock()
	defer mu.Unlock()

	if initialized == 0 {

		// Initialize a session that the SDK will use to load
		// credentials from the shared credentials file ~/.aws/credentials
		// and region from the shared configuration file ~/.aws/config.
		sess := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))

		// Create DynamoDB client
		svc := dynamodb.New(sess)
		dynamo = &clients.DynamoClient{svc: svc}
		atomic.StoreUint32(&dynamoInitialized, 1)
	}

	return dynamo
}

func GetS3Instance() clients.S3Client {

	if atomic.LoadUInt32(&s3ClientInitialized) == 1 {
		return s3client
	}

	mu.Lock()
	defer mu.Unlock()

	if initialized == 0 {
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String("eu-west-1")},
		)

		// Create S3 service client
		svc := s3.New(sess)

		s3client = &clients.S3Client{svc: svc}
		atomic.StoreUint32(&s3ClientInitialized, 1)
	}

	return s3client
}
