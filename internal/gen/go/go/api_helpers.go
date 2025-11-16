package openapi

import (
	"github.com/gin-gonic/gin"
)

// коды ошибок из openapi.yaml
const (
	errCodeTeamExists  = "TEAM_EXISTS"
	errCodePRExists    = "PR_EXISTS"
	errCodePRMerged    = "PR_MERGED"
	errCodeNotAssigned = "NOT_ASSIGNED"
	errCodeNoCandidate = "NO_CANDIDATE"
	errCodeNotFound    = "NOT_FOUND"
	errCodeInternal    = "INTERNAL_ERROR"
	errCodeBadRequest  = "BAD_REQUEST"
)

// error.response
type errorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type errorResponse struct {
	Error errorBody `json:"error"`
}

func writeError(c *gin.Context, status int, code, msg string) {
	c.JSON(status, errorResponse{
		Error: errorBody{
			Code:    code,
			Message: msg,
		},
	})
}
