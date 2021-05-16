package ddb

type FeedStat struct {
	FeedID        string `json:"FeedID"`
	LikesCount    int64  `json:"LikesCount"`
	CommentsCount int64  `json:"CommentsCount"`
	Version       int64  `json:"Version"`
}
