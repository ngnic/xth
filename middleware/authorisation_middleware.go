package middleware

import (
	"net/http"
	"strings"
	"xendit-takehome/github/repositories"

	"github.com/gin-gonic/gin"
)

func ApiKeyAuthorisation(userRepo repositories.UserRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeaderValue := ctx.Request.Header.Get("Authorization")
		if authHeaderValue == "" {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}
		dividerIndex := strings.Index(authHeaderValue, " ")
		if dividerIndex == -1 {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}
		prefix := authHeaderValue[:dividerIndex]
		apiKey := authHeaderValue[dividerIndex+1:]
		if prefix != "Bearer" || apiKey == "" {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		member, err := userRepo.GetUser(apiKey)
		if err != nil {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		ctx.Set("username", member.Username)
		ctx.Set("organisations", member.OrganisationNames)
		ctx.Next()
	}
}
