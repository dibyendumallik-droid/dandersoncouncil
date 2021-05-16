package handler

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/dandersoncouncil/covid_help/accessors"
	"github.com/dandersoncouncil/covid_help/clients"
	"github.com/dandersoncouncil/covid_help/datamodels/ddb"
	"github.com/dandersoncouncil/covid_help/handler/util"
)

type UserHandler struct {
	S3Accessor *accessors.S3Accessor
	DDbClient  *clients.DynamoClient
}

func (this *UserHandler) HandleRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	user := util.GetUser(req)

	switch req.HTTPMethod {
	case "GET":
		return this.get(user)
	case "POST":
		return this.create(user)
	default:
		return this.get(user)
	}
}

func (this *UserHandler) get(user *ddb.User) (events.APIGatewayProxyResponse, error) {
	userId := user.ID
	phoneNo := user.PhoneNo

	key := util.CreateProfilePicKey(userId, phoneNo)
	url := this.S3Accessor.CreatePreSignedGetUrl(util.BUCKET_NAME, key)

	user.ProfilePicLink = url

	response, err := json.Marshal(user)
	if err != nil {
		return util.ServerError(err)
	}

	// Return a response with a 200 OK status and the JSON book record
	// as the body.
	return util.CreateHttpRespose(string(response))
}

func (this *UserHandler) create(user *ddb.User) (events.APIGatewayProxyResponse, error) {
	userId := user.ID
	phoneNo := user.PhoneNo

	key := util.CreateProfilePicKey(userId, phoneNo)
	url := this.S3Accessor.CreatePreSignedGetUrl(util.BUCKET_NAME, key)

	user.ProfilePicLink = url

	this.DDbClient.CreateUser(user)

	return util.CreatePostSuccessRespose(userId)
}
