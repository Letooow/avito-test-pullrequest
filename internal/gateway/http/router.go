package http

import (
	handlers "avito-test/internal/gen/go/go"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc gin.HandlerFunc
}

func NewRouter(handleFunctions ApiHandleFunctions) *gin.Engine {
	return NewRouterWithGinEngine(gin.Default(), handleFunctions)
}

func getRoutes(handleFunctions ApiHandleFunctions) []Route {
	return []Route{
		{"PullRequestCreatePost", http.MethodPost, "/pullRequest/create", handleFunctions.PullRequestsAPI.PullRequestCreatePost},
		{"PullRequestMergePost", http.MethodPost, "/pullRequest/merge", handleFunctions.PullRequestsAPI.PullRequestMergePost},
		{"PullRequestReassignPost", http.MethodPost, "/pullRequest/reassign", handleFunctions.PullRequestsAPI.PullRequestReassignPost},
		{"TeamAddPost", http.MethodPost, "/team/add", handleFunctions.TeamsAPI.TeamAddPost},
		{"TeamGetGet", http.MethodGet, "/team/get", handleFunctions.TeamsAPI.TeamGetGet},
		{"UsersGetReviewGet", http.MethodGet, "/users/getReview", handleFunctions.UsersAPI.UsersGetReviewGet},
		{"UsersSetIsActivePost", http.MethodPost, "/users/setIsActive", handleFunctions.UsersAPI.UsersSetIsActivePost},
	}
}

func NewRouterWithGinEngine(router *gin.Engine, handleFunctions ApiHandleFunctions) *gin.Engine {
	for _, route := range getRoutes(handleFunctions) {
		if route.HandlerFunc == nil {
			route.HandlerFunc = handlers.DefaultHandleFunc
		}
		switch route.Method {
		case http.MethodGet:
			router.GET(route.Pattern, route.HandlerFunc)
		case http.MethodPost:
			router.POST(route.Pattern, route.HandlerFunc)
			// ...
		}
	}
	return router
}

type ApiHandleFunctions struct {
	PullRequestsAPI handlers.PullRequestsAPI
	TeamsAPI        handlers.TeamsAPI
	UsersAPI        handlers.UsersAPI
}
