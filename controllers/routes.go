package controllers

import (
	"xendit-takehome/github/repositories"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, orgRepo repositories.OrganisationRepository, authMiddleware gin.HandlerFunc) {
	orgController := OrganisationController{repository: orgRepo}
	orgRoutes := router.Group("/org")
	{
		orgRoutes.GET("/:org/comments/", orgController.getComments)
		orgRoutes.GET("/:org/members/", orgController.getMembers)
	}
	protectedOrgRoutes := router.Group("/org")
	protectedOrgRoutes.Use(authMiddleware)
	{

		protectedOrgRoutes.POST("/:org/comments/", orgController.postComments)
		protectedOrgRoutes.DELETE("/:org/comments/", orgController.deleteComments)
	}
}
