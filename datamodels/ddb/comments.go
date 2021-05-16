package ddb

type Comments struct {
	FeedID       string `json:"FeedID"`
	CommentID    string `json:CommentID`
	UserID       string `json:"UserID"`
	CreationTime int64  `json:"CreationTime"`
	Data         string `json:"Data"`
}
