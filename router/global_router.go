package router

import (
	"log"

	"github.com/dandersoncouncil/covid_help/handler"

	"github.com/aws/aws-lambda-go/events"
)

const (
	FEED_PATH               = "/Feed"
	FeedMediaPath           = "/Feed/Media"
	COMMENT_PATH            = "/Comment"
	USER_PATH               = "/User"
	PROFILE_PIC             = "/User/profile_pic"
	FeedStat                = "/FeedStat"
	CovidResourecPath       = "/covid_resouce"
	CovidResourceCategories = "/covid_resouce_categories"
)

type GlobalRouter struct {
	FeedHandler          *handler.FeedHandler
	PicHandler           *handler.PicHandler
	UserHandler          *handler.UserHandler
	CommentHandler       *handler.CommentHandler
	FeedStatHandler      *handler.FeedStatHandler
	FeedMediaHandler     *handler.FeedMediaHandler
	CovidResourcehandler *handler.CovidResourcehandler
	CategoryHandler      *handler.CategoryHandler
}

func (this *GlobalRouter) HandleRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	path := req.Path
	log.Printf("Path url: %s", path)

	if path == PROFILE_PIC {
		return this.PicHandler.HandleRequest(req)
	}

	if path == USER_PATH {
		return this.UserHandler.HandleRequest(req)
	}

	if path == COMMENT_PATH {
		return this.CommentHandler.HandleRequest(req)
	}

	if path == FeedStat {
		return this.FeedStatHandler.HandleRequest(req)
	}

	if path == FeedMediaPath {
		return this.FeedMediaHandler.HandleRequest(req)
	}

	if path == CovidResourecPath {
		return this.CovidResourcehandler.HandleRequest(req)
	}

	if path == CovidResourceCategories {
		return this.CategoryHandler.HandleRequest(req)
	}

	return this.FeedHandler.HandleRequest(req)
}
