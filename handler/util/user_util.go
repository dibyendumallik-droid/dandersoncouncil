package util

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/dandersoncouncil/covid_help/datamodels/ddb"
)

const (
	BUCKET_NAME        = "shareit-prod"
	PROFILE_PIC_FOLDER = "profile_pic/"
)

func CreateProfilePicKey(userId, phoneNo string) string {
	return fmt.Sprintf("%s%s%s", PROFILE_PIC_FOLDER, userId, phoneNo)
}

func GetUser(req events.APIGatewayProxyRequest) *ddb.User {
	context := req.RequestContext
	claims := context.Authorizer["claims"].(map[string]interface{})

	userId := fmt.Sprintf("%v", claims["cognito:username"])
	phoneNo := fmt.Sprintf("%v", claims["phone_number"])
	email := fmt.Sprintf("%s", claims["email"])
	name := fmt.Sprintf("%v", claims["name"])

	return &ddb.User{userId, name, email, phoneNo, "", true}
}
