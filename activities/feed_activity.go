package activities

import (
	"net/http"

	"../clients"
)

type FeedActivity struct {
	dynamoClient *clients.DynamoClient
	s3Client     *clients.S3Client
}

func (activity *FeedActivity) ListFeeds(w http.ResponseWriter, r *http.Request) {
	// feeds := activity.dynamoClient.ListFeeds()

	// feedResponses := []*responses.GetFeedResponse{}

	// for _, feed := range feeds {
	// 	id := feed.ID
	// 	images := feed.MediaList

	// 	signedUrls := []string{}
	// 	for _, img := range images {
	// 		var timeDuration time.Duration = 10000
	// 		signedURL := activity.s3Client.GetSignedURL(img.Bucket, img.Key, timeDuration)
	// 		signedUrls = append(signedUrls, signedURL)
	// 	}
	// 	feedResponse := &responses.GetFeedResponse{ID: id, URLs: signedUrls}
	// 	feedResponses = append(feedResponses, feedResponse)
	// }

	// listResponse := responses.ListFeedResponse{Feeds: feedResponses}
}
