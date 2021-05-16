package ddb

type Media struct {
	Type string `json:"Type"`
	URL  string `json:"url"`
}

const (
	MediaTypeVideo = "Video"
	MediaTypeImage = "Image"
)
