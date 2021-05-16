package accessors

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/dandersoncouncil/covid_help/datamodels/ddb"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	elastic "github.com/olivere/elastic/v7"
)

// Elasticsearch is an ES Client which will perform Elasticsearch Updates for Dynamo Items
type Elasticsearch struct {
	*elastic.Client
}

// Update takes a reference to adstream.Details object;
// which is used to figure out which Elasticsearch Index to update;
// And an item map[string]events.DynamoDBAttributeValue which will be turned into JSON
// then indexed into Elasticsearch
func (e *Elasticsearch) Update(d *Details, item map[string]events.DynamoDBAttributeValue) (*elastic.IndexResponse, error) {
	//log.Printf("Trying to create a new index for covid")
	//e.CreateIndexFoeCR(d)

	tmp := eventStreamToMap(item)
	var i interface{}
	if err := dynamodbattribute.UnmarshalMap(tmp, &i); err != nil {
		return nil, err

	}
	resp, err := e.Index().
		Id(d.docID(item)).
		Index(d.index()).
		BodyJson(i).
		Type("_doc").
		Do(context.Background())

	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return resp, nil
}

func (e *Elasticsearch) QueryNearByCovidResources(d *Details, location *ddb.Location, offset, limit int, category string) []*ddb.CovidResource {
	geoDistanceSorter := elastic.NewGeoDistanceSort("Loc").
		Point(location.Lat, location.Lon).
		Unit("km").
		GeoDistance("plane").
		Asc()

	//matchAllQuery := elastic.NewMatchAllQuery()
	termQuery := elastic.NewTermQuery("Category", category)

	//searchResult, err:=e.Search().Index(d.index()).Query(elastic.NewMatchAllQuery()).Type("_doc")
	searchResult, err := e.Search().
		Index(d.index()).
		Query(termQuery).
		//Size(limit).
		//From(offset).
		SortBy(geoDistanceSorter).
		Do(context.Background())

	if err != nil {
		v, _ := geoDistanceSorter.Source()
		log.Printf("Error in QueryNearByCovidResources: %v %v", err, v)
		return []*ddb.CovidResource{}
	}
	log.Printf("searchResult:  %v", searchResult)

	var feedList = []*ddb.CovidResource{}

	var ttyp ddb.CovidResource

	for _, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
		t := item.(ddb.CovidResource)
		feedList = append(feedList, &t)
	}
	return feedList
}

func (e *Elasticsearch) QueryNearByFeedElement(d *Details, location *ddb.Location) []*ddb.Feed {
	geoDistanceSorter := elastic.NewGeoDistanceSort("location").
		Point(location.Lat, location.Lon).
		Unit("km").
		GeoDistance("plane").
		Asc()

	matchAllQuery := elastic.NewMatchAllQuery()
	searchResult, err := e.Search().
		Index(d.index()).
		Query(matchAllQuery).
		Size(20).
		SortBy(geoDistanceSorter).
		Do(context.Background())

	if err != nil {
		v, _ := geoDistanceSorter.Source()
		log.Printf("Error in QueryNearByFeedElement: %v %v", err, v)
		return []*ddb.Feed{}
	}
	log.Printf("searchResult:  %v", searchResult)

	var feedList = []*ddb.Feed{}

	var ttyp ddb.Feed

	for _, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
		t := item.(ddb.Feed)
		feedList = append(feedList, &t)
	}
	return feedList
}

// create index for covid
// use this  ethod for one time creation of elastic search cluster
func (e *Elasticsearch) CreateIndexFoeCR(d *Details) {
	/*
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
	*/

	/*
			ID           string   `json:"ID"`
		Data         string   `json:"Data"`
		CreatedBy    string   `json:"CreatedBy"`
		CreatedByID  string   `json:"CreatedByID"`
		CreationTime int64    `json:"CreationTime"`
		Category     string   `json:"Category"`
		MediaList    []Media  `json:"MediaList"`
		Location     Location `json:"location"`
	*/

	mapping := `
	{
		"mappings":{
				"properties":{
					"Data": {
						"type":"text"
					},
					"ID": {
						"type":"keyword"
					},
					"CreatedBy": {
						"type":"keyword"
					},
					"CreatedByID": {
						"type":"keyword"
					},
					"MediaList": {
						"type":"keyword"
					},
					"Name": {
						"type":"keyword"
					},
					"Category": {
						"type":"keyword"
					},
					"AddrLine": {
						"type":"text"
					},
					"City": {
						"type":"keyword"
					},
					"State": {
						"type":"keyword"
					},
					"PhoenNo":{
						"type":"keyword"
					},
					"IsVerfied":{
						"type":"boolean"
					},
					"LastVerifiedTime":{
						"type":"keyword"
					},
					"Remarks":{
						"type":"keyword"
					},
					"ConfidenceScore":{
						"type":"integer"
					},
					"Loc":{
						"type":"geo_point"
					}
				}
		}
	}`

	e.DeleteIndex(d.index()).Do(context.Background())

	createIndex, err := e.CreateIndex(d.index()).Body(mapping).Do(context.Background())
	log.Printf("createIndex: %v", createIndex)

	if err != nil {
		log.Printf("Error: %v ", err)
	}
}

// func (e *Elasticsearch) CreateIndex(d *Details) {
// 	/*
// 			ID           string   `json:"ID"`
// 		Data         string   `json:"Data"`
// 		CreatedBy    string   `json:"CreatedBy"`
// 		CreatedByID  string   `json:"CreatedByID"`
// 		CreationTime int64    `json:"CreationTime"`
// 		Category     string   `json:"Category"`
// 		MediaList    []Media  `json:"MediaList"`
// 		Location     Location `json:"location"`
// 	*/

// 	e.DeleteIndex(d.index()).Do(context.Background())
// 	// createIndex, err := e.CreateIndex(d.index()).Body(mapping).Do(context.Background())
// 	// log.Printf("createIndex: %v", createIndex)

// 	// if err != nil {
// 	// 	log.Printf("Error: %v ", err)
// 	// }

// }

func (d *Details) docID(item map[string]events.DynamoDBAttributeValue) (id string) {
	if d != nil {
		if d.RangeKey != "" {
			id = fmt.Sprintf("%s-%s", item[d.HashKey].String(), item[d.RangeKey].String())
		} else {
			id = item[d.HashKey].String()
		}
	}
	return id
}

func (d *Details) index() string {
	return fmt.Sprintf("%sindex", strings.ToLower(d.TableName))
}

// ugly hack because the types
// events.DynamoDBAttributeValue != *dynamodb.AttributeValue
func eventStreamToMap(attribute interface{}) map[string]*dynamodb.AttributeValue {
	// Map to be returned
	m := make(map[string]*dynamodb.AttributeValue)

	tmp := make(map[string]events.DynamoDBAttributeValue)

	switch t := attribute.(type) {
	case map[string]events.DynamoDBAttributeValue:
		tmp = t
	case events.DynamoDBAttributeValue:
		tmp = t.Map()
	}

	for k, v := range tmp {
		switch v.DataType() {
		case events.DataTypeString:
			s := v.String()
			m[k] = &dynamodb.AttributeValue{
				S: &s,
			}
		case events.DataTypeBoolean:
			b := v.Boolean()
			m[k] = &dynamodb.AttributeValue{
				BOOL: &b,
			}
		case events.DataTypeMap:
			m[k] = &dynamodb.AttributeValue{
				M: eventStreamToMap(v),
			}
		case events.DataTypeNumber:
			n := v.Number()
			m[k] = &dynamodb.AttributeValue{
				N: &n,
			}
		case events.DataTypeList:
			m[k] = &dynamodb.AttributeValue{
				L: eventStreamToList(v),
			}
		}
	}
	return m
}

// ugly hack because the types
// events.DynamoDBAttributeValue != *dynamodb.AttributeValue
func eventStreamToList(attribute interface{}) []*dynamodb.AttributeValue {
	// List to be returned
	l := make([]*dynamodb.AttributeValue, 0)

	var tmp []events.DynamoDBAttributeValue

	switch t := attribute.(type) {
	case []events.DynamoDBAttributeValue:
		tmp = t
	case events.DynamoDBAttributeValue:
		tmp = t.List()
	}

	for _, v := range tmp {
		switch v.DataType() {
		case events.DataTypeString:
			s := v.String()
			l = append(l, &dynamodb.AttributeValue{
				S: &s,
			})
		case events.DataTypeBoolean:
			b := v.Boolean()
			l = append(l, &dynamodb.AttributeValue{
				BOOL: &b,
			})
		case events.DataTypeMap:
			l = append(l, &dynamodb.AttributeValue{
				M: eventStreamToMap(v),
			})
		case events.DataTypeNumber:
			n := v.Number()
			l = append(l, &dynamodb.AttributeValue{
				N: &n,
			})
		case events.DataTypeList:
			l = append(l, &dynamodb.AttributeValue{
				L: eventStreamToList(v),
			})
		}
	}
	return l
}
