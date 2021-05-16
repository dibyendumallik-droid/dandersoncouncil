package ddb

type CovidResource struct {
	ResourceId       string   `json:"ID"` //Primary identifier concatanation of name + phone number
	Name             string   `json:"Name"`
	Category         string   `json:"Category"`
	Loc              Location `json:"Loc"`
	AddrLine         string   `json:"AddrLine"`
	City             string   `json:"City"`
	State            string   `json:"State"`
	PhoenNo          string   `json:"PhoenNo"`
	IsVerfied        bool     `json:"IsVerfied"`
	LastVerifiedTime int64    `json:"LastVerifiedTime"`
	Remarks          string   `json:"Remarks"`
	ConfidenceScore  int64    `json:"ConfidenceScore"` // For internal use only
}
