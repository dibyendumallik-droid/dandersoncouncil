package ddb

type Offer struct {
	ResourceId         string   `json:"ResourceId"` // hash key
	OfferId            string   `json:"OfferId"`    //Part of sort key
	MerchantId         string   `json:"MerchantId"` // Part of sort key
	Price              float64  `json:"Price"`
	Qty                int64    `json:"Qty"`
	Unit               string   `json:"Unit"`
	AvailableUnitCount int64    `json:"AvailableUnit"`
	Loc                Location `json:"Loc"`
	ServingRad         float64  `json:"ServingRad"`
}
