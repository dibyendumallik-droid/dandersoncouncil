package handler

import (
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/dandersoncouncil/covid_help/accessors"
	"github.com/dandersoncouncil/covid_help/handler/util"
)

type PicHandler struct {
	S3Accessor *accessors.S3Accessor
}

func (this *PicHandler) HandleRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// environment variables
	context := req.RequestContext
	claims := context.Authorizer["claims"].(map[string]interface{})

	userId := fmt.Sprintf("%v", claims["cognito:username"])
	phoneNo := fmt.Sprintf("%v", claims["phone_number"])

	log.Printf("userId: %s", userId)
	log.Printf("phoneNo: %s", phoneNo)

	//key := fmt.Sprintf("%s%s%s", PROFILE_PIC_FOLDER, userId, phoneNo)
	key := util.CreateProfilePicKey(userId, phoneNo)
	var url string
	switch req.HTTPMethod {
	case "GET":
		url = this.S3Accessor.CreatePreSignedGetUrl(util.BUCKET_NAME, key)
	case "POST":
		url = this.S3Accessor.CreatePreSignedPutUrl(util.BUCKET_NAME, key)
	default:
		url = this.S3Accessor.CreatePreSignedPutUrl(util.BUCKET_NAME, key)
	}
	return util.CreateHttpRespose(url)
}
