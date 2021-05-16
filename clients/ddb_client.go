package clients

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/dandersoncouncil/covid_help/datamodels/ddb"
)

const (
	FeedTableName                 = "Feed"
	UserTableName                 = "User"
	CommentsTableName             = "Comment"
	CommentsTableHashKey          = "FeedID"
	LikeTableName                 = "Like"
	LikeCountTableName            = "LikeCountTableName"
	FeedStatTableName             = "FeedStat"
	FeedStatTableNameHashKey      = "FeedID"
	CovidRelatedResourceTableName = "CovidResources"
	CovidRelatedResourceHashKey   = "ID"
	OfferTableName                = "Offers"
)

type DynamoClient struct {
	SVC *dynamodb.DynamoDB
}

func (client *DynamoClient) GetOfferData(order ddb.Order) (*ddb.Offer, error) {
	resourceId := order.ResourceId
	offerId := order.OfferId

	result, err := client.SVC.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(OfferTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"ResourecId": {
				S: aws.String(resourceId),
			},
			"OfferId": {
				S: aws.String(offerId),
			},
		},
	})

	if err != nil {
		return nil, err
	}

	item := new(ddb.Offer)
	err = dynamodbattribute.UnmarshalMap(result.Item, item)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (client *DynamoClient) RemoveCovidResource(resourceID string) error {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			CovidRelatedResourceHashKey: {
				S: aws.String(resourceID),
			},
		},
		TableName: aws.String(CovidRelatedResourceTableName),
	}

	_, err := client.SVC.DeleteItem(input)

	return err
}

func (c *DynamoClient) CreateCovidResource(data *ddb.CovidResource) bool {
	id := fmt.Sprintf("%v-%v", data.Name, data.PhoenNo)
	data.ResourceId = id
	av, err := dynamodbattribute.MarshalMap(data)
	if err != nil {
		log.Println("Got error marshalling map:")
		log.Println(err.Error())
		return false
	}

	// Create item in table Movies
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(CovidRelatedResourceTableName),
	}

	_, err = c.SVC.PutItem(input)

	if err != nil {
		log.Printf("%v", err)
		return false
	}
	return true
}

func (this *DynamoClient) CreateUser(data *ddb.User) bool {
	av, err := dynamodbattribute.MarshalMap(data)
	if err != nil {
		log.Println("Got error marshalling map:")
		log.Println(err.Error())
		return false
	}

	// Create item in table Movies
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(UserTableName),
	}

	_, err = this.SVC.PutItem(input)

	if err != nil {
		log.Printf("%v", err)
		return false
	}
	return true
}

func (this *DynamoClient) CreateFeed(data *ddb.Feed) bool {
	av, err := dynamodbattribute.MarshalMap(data)
	if err != nil {
		fmt.Println("Got error marshalling map:")
		fmt.Println(err.Error())
		return false
	}

	// Create item in table Movies
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(FeedTableName),
	}

	_, err = this.SVC.PutItem(input)

	if err != nil {
		log.Printf("%v", err)
		return false
	}
	return true
}

func (this *DynamoClient) RemoveFeed(feedID string) error {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"FeedID": {
				S: aws.String(feedID),
			},
		},
		TableName: aws.String(FeedTableName),
	}

	_, err := this.SVC.DeleteItem(input)

	return err
}

func (client *DynamoClient) ListComments(feedId string) []*ddb.Comments {
	comments := []*ddb.Comments{}
	// Build the query input parameters

	queryInput := &dynamodb.QueryInput{
		KeyConditionExpression: aws.String("FeedID = :pk"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pk": {
				S: aws.String(feedId),
			},
		},
		TableName: aws.String(CommentsTableName),
	}

	result, err := client.SVC.Query(queryInput)

	if err != nil {
		log.Printf("Error while getting feed data: %v", err.Error())
	}

	for _, i := range result.Items {
		comment := ddb.Comments{}

		err = dynamodbattribute.UnmarshalMap(i, &comment)

		if err != nil {
			log.Printf("Error while getting feed data: %v", err.Error())
		}

		comments = append(comments, &comment)
	}

	return comments
}

func (client *DynamoClient) CreateComment(feedID, commentID, userID, data string) bool {
	now := time.Now()
	currentEpochTimeInMillis := now.UnixNano() / 1000000
	comments := &ddb.Comments{feedID, commentID, userID, currentEpochTimeInMillis, data}
	av, err := dynamodbattribute.MarshalMap(comments)
	if err != nil {
		log.Println("Got error marshalling map:")
		log.Println(err.Error())
		return false
	}

	// Create item in table Movies
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(CommentsTableName),
	}

	_, err = client.SVC.PutItem(input)

	if err != nil {
		log.Printf("%v", err)
		return false
	}
	return true
}

func (this *DynamoClient) RemoveComment(feedID, commentID string) error {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"FeedID": {
				S: aws.String(feedID),
			},
			"CommentID": {
				S: aws.String(commentID),
			},
		},
		TableName: aws.String(CommentsTableName),
	}

	_, err := this.SVC.DeleteItem(input)

	return err
}

func (client *DynamoClient) UpdateLikeCount(likesCount int16, feedId string) error {
	updateExpression := "SET LikesCount = LikesCount + :likesCount"
	update := &dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			FeedStatTableNameHashKey: &dynamodb.AttributeValue{
				S: aws.String(feedId),
			},
		},
		TableName: aws.String(FeedStatTableName),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":likesCount": &dynamodb.AttributeValue{
				N: aws.String(fmt.Sprintf("%v", likesCount)),
			},
		},
		UpdateExpression: aws.String(updateExpression),
	}
	_, err := client.SVC.UpdateItem(update)

	if err != nil {
		log.Printf("Erro while updating likes %v", err)
		return err
	}
	return nil
}

func (client *DynamoClient) UpdateCommentsCount(commentsCount int16, feedId string) error {
	updateExpression := "SET CommentsCount = if_not_exists(CommentsCount, :zero) + :commentsCount"
	update := &dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			FeedStatTableNameHashKey: &dynamodb.AttributeValue{
				S: aws.String(feedId),
			},
		},
		TableName: aws.String(FeedStatTableName),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":commentsCount": &dynamodb.AttributeValue{
				N: aws.String(fmt.Sprintf("%v", commentsCount)),
			},
			":zero": &dynamodb.AttributeValue{
				N: aws.String(fmt.Sprintf("%v", 0)),
			},
		},
		UpdateExpression: aws.String(updateExpression),
	}
	_, err := client.SVC.UpdateItem(update)

	if err != nil {
		log.Printf("Erro while updating comments %v", err)
		return err
	}
	return nil
}

func (this *DynamoClient) CreateLike(data *ddb.Like) error {
	av, err := dynamodbattribute.MarshalMap(data)
	if err != nil {
		return err
	}

	// Create item in table Movies
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(LikeTableName),
	}

	_, err = this.SVC.PutItem(input)

	if err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

func (this *DynamoClient) RemoveLike(data *ddb.Like) error {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"FeedID": {
				S: aws.String(data.FeedID),
			},
			"UserID": {
				S: aws.String(data.UserID),
			},
		},
		TableName: aws.String(LikeTableName),
	}

	_, err := this.SVC.DeleteItem(input)

	return err
}

func (client *DynamoClient) GetFeedStat(feedID string) (*ddb.FeedStat, error) {
	result, err := client.SVC.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(FeedStatTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"FeedID": {
				S: aws.String(feedID),
			},
		},
	})

	if err != nil {
		return nil, err
	}

	item := new(ddb.FeedStat)

	err = dynamodbattribute.UnmarshalMap(result.Item, item)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (client *DynamoClient) PathFeed(feedID string, mediaList []string) error {
	result, err := client.SVC.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(FeedTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(feedID),
			},
		},
	})

	if err != nil {
		return err
	}

	storedFeed := new(ddb.Feed)
	err = dynamodbattribute.UnmarshalMap(result.Item, storedFeed)
	storedMediaList := storedFeed.MediaList
	mergedList := append(storedMediaList, mediaList...)
	storedFeed.MediaList = mergedList

	av, err := dynamodbattribute.MarshalMap(storedFeed)
	if err != nil {
		log.Println("Got error marshalling map:")
		log.Println(err.Error())
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(FeedTableName),
	}

	_, err = client.SVC.PutItem(input)

	return err
}
