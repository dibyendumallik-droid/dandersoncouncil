package ddb

type CovidResource struct {
	ResourceId       string   `json:"ID"` //Primary identifier concatanation of name + phone number, it is a foreign key in offer table
	Name             string   `json:"Name"`
	Category         string   `json:"Category"` // For example it can be oxygen concentrator
	Loc              Location `json:"Loc"`
	AddrLine         string   `json:"AddrLine"`
	City             string   `json:"City"`
	State            string   `json:"State"`
	PhoenNo          string   `json:"PhoenNo"`
	IsVerfied        bool     `json:"IsVerfied"`
	LastVerifiedTime int64    `json:"LastVerifiedTime"`
	Remarks          string   `json:"Remarks"`
	ConfidenceScore  int64    `json:"ConfidenceScore"` // For internal use only

	TotalResourceCount    int64 `json:"TotalResourceCount"`
	OccupiedResourceCount int64 `json:"OccupiedResourceCount"`
}
