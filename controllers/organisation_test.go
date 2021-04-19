package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"xendit-takehome/github/controllers"
	"xendit-takehome/github/entities"
	mocks "xendit-takehome/github/mocks/repositories"
	"xendit-takehome/github/responses"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetCommentsRoute(t *testing.T) {
	orgRepo := mocks.OrganisationRepository{}
	comments := []entities.Comment{
		{
			ID:             1,
			Comment:        "ABC",
			CreatedAt:      time.Now(),
			CreatedBy:      1,
			OrganisationId: 1,
		},
		{
			ID:             2,
			Comment:        "DEF",
			CreatedAt:      time.Now(),
			CreatedBy:      1,
			OrganisationId: 1,
		},
	}
	orgRepo.On("GetComments", "A").Return(comments, nil)
	gin.SetMode(gin.TestMode)

	responseRecorder := httptest.NewRecorder()
	_, router := gin.CreateTestContext(responseRecorder)
	controllers.SetupRoutes(router, &orgRepo, func(ctx *gin.Context) {
		// Just a stub
		ctx.Next()
	})

	//given
	request, _ := http.NewRequest(http.MethodGet, "/org/A/comments/", nil)

	//when
	router.ServeHTTP(responseRecorder, request)

	//then
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	expectedResponse := responses.CreateGetCommentsResponse(comments)
	expectedResponseBytes, _ := json.Marshal(expectedResponse)
	assert.JSONEq(t, string(expectedResponseBytes), responseRecorder.Body.String())
	orgRepo.AssertExpectations(t)
}

func TestGetCommentsRouteNoResult(t *testing.T) {
	orgRepo := mocks.OrganisationRepository{}
	comments := []entities.Comment{}
	orgRepo.On("GetComments", "A").Return(comments, nil)
	gin.SetMode(gin.TestMode)

	responseRecorder := httptest.NewRecorder()
	_, router := gin.CreateTestContext(responseRecorder)
	controllers.SetupRoutes(router, &orgRepo, func(ctx *gin.Context) {
		// Just a stub
		ctx.Next()
	})

	//given
	request, _ := http.NewRequest(http.MethodGet, "/org/A/comments/", nil)

	//when
	router.ServeHTTP(responseRecorder, request)

	//then
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	expectedResponse := responses.CreateGetCommentsResponse(comments)
	expectedResponseBytes, _ := json.Marshal(expectedResponse)
	assert.JSONEq(t, string(expectedResponseBytes), responseRecorder.Body.String())
	orgRepo.AssertExpectations(t)
}

func TestGetMembersRoute(t *testing.T) {
	orgRepo := mocks.OrganisationRepository{}
	members := []entities.Member{
		{
			ID:           1,
			Username:     "A",
			Avatar:       "A",
			CreatedAt:    time.Now(),
			PasswordHash: "A",
			Following:    1,
			FollowedBy:   1,
		},
		{
			ID:           2,
			Username:     "B",
			Avatar:       "B",
			CreatedAt:    time.Now(),
			PasswordHash: "A",
			Following:    1,
			FollowedBy:   1,
		},
	}
	orgRepo.On("GetMembers", "A").Return(members, nil)
	gin.SetMode(gin.TestMode)

	responseRecorder := httptest.NewRecorder()
	_, router := gin.CreateTestContext(responseRecorder)
	controllers.SetupRoutes(router, &orgRepo, func(ctx *gin.Context) {
		// Just a stub
		ctx.Next()
	})

	//given
	request, _ := http.NewRequest(http.MethodGet, "/org/A/members/", nil)

	//when
	router.ServeHTTP(responseRecorder, request)

	//then
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	expectedResponse := responses.CreateGetMembersResponse(members)
	expectedResponseBytes, _ := json.Marshal(expectedResponse)
	assert.JSONEq(t, string(expectedResponseBytes), responseRecorder.Body.String())
	orgRepo.AssertExpectations(t)
}

func TestGetMembersRouteNoResult(t *testing.T) {
	orgRepo := mocks.OrganisationRepository{}
	members := []entities.Member{}
	orgRepo.On("GetMembers", "A").Return(members, nil)
	gin.SetMode(gin.TestMode)

	responseRecorder := httptest.NewRecorder()
	_, router := gin.CreateTestContext(responseRecorder)
	controllers.SetupRoutes(router, &orgRepo, func(ctx *gin.Context) {
		// Just a stub
		ctx.Next()
	})

	//given
	request, _ := http.NewRequest(http.MethodGet, "/org/A/members/", nil)

	//when
	router.ServeHTTP(responseRecorder, request)

	//then
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	expectedResponse := responses.CreateGetMembersResponse(members)
	expectedResponseBytes, _ := json.Marshal(expectedResponse)
	assert.JSONEq(t, string(expectedResponseBytes), responseRecorder.Body.String())
	orgRepo.AssertExpectations(t)
}

func TestPostCommentsRoute(t *testing.T) {
	org := "orgA"
	username := "userA"
	comment := "ABC"

	orgRepo := mocks.OrganisationRepository{}
	orgRepo.On("AddComment", org, username, comment).Return(nil)
	gin.SetMode(gin.TestMode)

	responseRecorder := httptest.NewRecorder()
	_, router := gin.CreateTestContext(responseRecorder)
	controllers.SetupRoutes(router, &orgRepo, func(ctx *gin.Context) {
		ctx.Set("username", username)
		ctx.Set("organisations", []string{org})
		ctx.Next()
	})

	requestBody := map[string]interface{}{
		"comment": comment,
	}
	encodedBody, _ := json.Marshal(requestBody)

	//given
	request, _ := http.NewRequest(http.MethodPost, "/org/orgA/comments/", bytes.NewReader(encodedBody))
	request.Header.Set("Content-Type", "application/json")

	//when
	router.ServeHTTP(responseRecorder, request)

	//then
	assert.Equal(t, http.StatusCreated, responseRecorder.Code)
	orgRepo.AssertExpectations(t)
}
func TestPostCommentsRouteReturn400OnValidationErrors(t *testing.T) {
	userName := "orgA"

	orgRepo := mocks.OrganisationRepository{}
	gin.SetMode(gin.TestMode)

	t.Run("when request body is empty", func(t *testing.T) {
		responseRecorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(responseRecorder)
		controllers.SetupRoutes(router, &orgRepo, func(ctx *gin.Context) {
			ctx.Set("username", userName)
			ctx.Next()
		})

		requestBody := map[string]interface{}{}
		encodedBody, _ := json.Marshal(requestBody)

		//given
		request, _ := http.NewRequest(http.MethodPost, "/org/A/comments/", bytes.NewReader(encodedBody))
		request.Header.Set("Content-Type", "application/json")

		//when
		router.ServeHTTP(responseRecorder, request)

		//then
		assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
		orgRepo.AssertExpectations(t)
	})

	t.Run("when comment is empty", func(t *testing.T) {
		responseRecorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(responseRecorder)
		controllers.SetupRoutes(router, &orgRepo, func(ctx *gin.Context) {
			ctx.Set("username", userName)
			ctx.Next()
		})

		requestBody := map[string]interface{}{
			"comment": "",
		}
		encodedBody, _ := json.Marshal(requestBody)

		//given
		request, _ := http.NewRequest(http.MethodPost, "/org/A/comments/", bytes.NewReader(encodedBody))
		request.Header.Set("Content-Type", "application/json")

		//when
		router.ServeHTTP(responseRecorder, request)

		//then
		assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
		orgRepo.AssertExpectations(t)
	})

	t.Run("when repository raises an error", func(t *testing.T) {
		org := "A"
		comment := "ABC"
		orgRepo.On("AddComment", org, userName, comment).Return(errors.New("db error"))
		responseRecorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(responseRecorder)
		controllers.SetupRoutes(router, &orgRepo, func(ctx *gin.Context) {
			ctx.Set("username", userName)
			ctx.Next()
		})

		requestBody := map[string]interface{}{
			"comment": comment,
		}
		encodedBody, _ := json.Marshal(requestBody)

		//given
		request, _ := http.NewRequest(http.MethodPost, "/org/A/comments/", bytes.NewReader(encodedBody))
		request.Header.Set("Content-Type", "application/json")

		//when
		router.ServeHTTP(responseRecorder, request)

		//then
		assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
		orgRepo.AssertExpectations(t)
	})
}

func TestDeleteCommentsRoute(t *testing.T) {
	org := "orgA"
	username := "userA"

	orgRepo := mocks.OrganisationRepository{}
	orgRepo.On("DeleteComments", org).Return(nil)
	gin.SetMode(gin.TestMode)

	responseRecorder := httptest.NewRecorder()
	_, router := gin.CreateTestContext(responseRecorder)
	controllers.SetupRoutes(router, &orgRepo, func(ctx *gin.Context) {
		ctx.Set("username", username)
		ctx.Set("organisations", []string{org})
		ctx.Next()
	})

	//given
	request, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/org/%s/comments/", org), nil)

	//when
	router.ServeHTTP(responseRecorder, request)

	//then
	assert.Equal(t, http.StatusNoContent, responseRecorder.Code)
	orgRepo.AssertExpectations(t)
}

func TestDeleteCommentsRouteErrorHandling(t *testing.T) {
	org := "orgA"
	orgRepo := mocks.OrganisationRepository{}
	gin.SetMode(gin.TestMode)

	t.Run("when middleware is misconfigured", func(t *testing.T) {
		responseRecorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(responseRecorder)
		controllers.SetupRoutes(router, &orgRepo, func(ctx *gin.Context) {
			ctx.Next()
		})

		//given
		request, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/org/%s/comments/", org), nil)

		//when
		router.ServeHTTP(responseRecorder, request)

		//then
		assert.Equal(t, http.StatusInternalServerError, responseRecorder.Code)
		orgRepo.AssertExpectations(t)
	})

	t.Run("when member isn't part of the organisation", func(t *testing.T) {
		username := "A"
		responseRecorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(responseRecorder)
		controllers.SetupRoutes(router, &orgRepo, func(ctx *gin.Context) {
			ctx.Set("username", username)
			ctx.Set("organisations", []string{})
			ctx.Next()
		})

		//given
		request, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/org/%s/comments/", org), nil)

		//when
		router.ServeHTTP(responseRecorder, request)

		//then
		assert.Equal(t, http.StatusForbidden, responseRecorder.Code)
		orgRepo.AssertExpectations(t)
	})

	t.Run("when db raises an error", func(t *testing.T) {
		username := "A"
		responseRecorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(responseRecorder)
		controllers.SetupRoutes(router, &orgRepo, func(ctx *gin.Context) {
			ctx.Set("username", username)
			ctx.Set("organisations", []string{org})
			ctx.Next()
		})

		orgRepo.On("DeleteComments", org).Return(errors.New("db error"))

		//given
		request, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/org/%s/comments/", org), nil)

		//when
		router.ServeHTTP(responseRecorder, request)

		//then
		assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
		orgRepo.AssertExpectations(t)
	})
}
