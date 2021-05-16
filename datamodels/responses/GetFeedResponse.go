package responses

type GetFeedResponse struct {
	ID   string   `json:"ID"`
	URLs []string `json:"URLs"`
}
