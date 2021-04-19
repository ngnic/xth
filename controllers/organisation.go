package controllers

import (
	"errors"
	"net/http"
	"xendit-takehome/github/repositories"
	"xendit-takehome/github/requests"
	"xendit-takehome/github/responses"

	"github.com/gin-gonic/gin"
)

type OrganisationController struct {
	repository repositories.OrganisationRepository
}

func (controller *OrganisationController) getComments(ctx *gin.Context) {
	org := ctx.Param("org")
	comments, _ := controller.repository.GetComments(org)
	ctx.JSON(http.StatusOK, responses.CreateGetCommentsResponse(comments))
}

func (controller *OrganisationController) getMembers(ctx *gin.Context) {
	org := ctx.Param("org")
	members, _ := controller.repository.GetMembers(org)
	ctx.JSON(http.StatusOK, responses.CreateGetMembersResponse(members))
}

func (controller *OrganisationController) postComments(ctx *gin.Context) {
	var request requests.PostCommentRequest
	username, _ := ctx.Get("username")
	requestedOrg := ctx.Param("org")
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := controller.repository.AddComment(requestedOrg, username.(string), request.Comment); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusCreated)
}

func (controller *OrganisationController) deleteComments(ctx *gin.Context) {
	requestedOrg, err := checkUserInOrg(ctx)
	if err != nil && errors.Is(err, ErrMisconfiguredAuthorisation) {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected error"})
		return
	}

	if err != nil && errors.Is(err, ErrMismatchedOrg) {
		ctx.Status(http.StatusForbidden)
		return
	}
	if err := controller.repository.DeleteComments(requestedOrg); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}
