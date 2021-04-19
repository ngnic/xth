package controllers

import (
	"errors"

	"github.com/gin-gonic/gin"
)

var (
	ErrMismatchedOrg              = errors.New("user does not have permission")
	ErrMisconfiguredAuthorisation = errors.New("unexpected error has occured")
)

func checkUserInOrg(ctx *gin.Context) (string, error) {
	requestedOrg := ctx.Param("org")
	userOrganisations, exists := ctx.Get("organisations")
	if !exists {
		return "", ErrMisconfiguredAuthorisation
	}

	hasMatchingOrg := false
	for _, org := range userOrganisations.([]string) {
		if org == requestedOrg {
			hasMatchingOrg = true
			break
		}
	}

	if !hasMatchingOrg {
		return "", ErrMismatchedOrg
	}

	return requestedOrg, nil
}
