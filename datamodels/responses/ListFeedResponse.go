package responses

import "../ddb"

type ListFeedResponse struct {
	Feeds []*ddb.Feed `json:"feed"`
}
