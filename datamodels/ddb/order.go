package ddb

type Order struct {
	OrderId      string  `json:"OrderId"`
	ResourceId   string  `json:"ResourceId"`
	OfferId      string  `json:"OfferId"`
	MerchantId   string  `json:"MerchantId"`
	UnitCount    int64   `json:"UnitCount"`
	Price        float64 `json:"Price"`
	DeliveryAddr string  `json:"DeliveryAddr"`
	CustomerId   string  `json:"CustomerId"`
}
