package ddb

type Feed struct {
	ID             string   `json:"ID"`
	Data           string   `json:"Data"`
	CreatedBy      string   `json:"CreatedBy"`
	CreatedByID    string   `json:"CreatedByID"`
	CreationTime   int64    `json:"CreationTime"`
	Category       string   `json:"Category"`
	MediaList      []string `json:"MediaList"`
	Location       Location `json:"Loc"`
	OfferStartTime int64    `json:"OfferStartTime"`
	OfferEndTime   int64    `json:"OfferEndTime"`
	IsExpired      bool     `json:"IsExpired"`
}
