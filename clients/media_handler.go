package clients

import (
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Client struct {
	svc *s3.S3
}

func (s *S3Client) GetSignedURL(bucket string, key string, timeInMinutes time.Duration) string {
	// Initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	req, _ := s.svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	urlStr, err := req.Presign(timeInMinutes * time.Minute)

	if err != nil {
		log.Println("Failed to sign request", err)
	}

	return urlStr
}
