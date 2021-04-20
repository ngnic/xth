package middleware_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"xendit-takehome/github/entities"
	"xendit-takehome/github/middleware"
	mocks "xendit-takehome/github/mocks/repositories"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestAuthorisationMiddleware(t *testing.T) {
	mockUserRepo := &mocks.UserRepository{}
	mockUserRepo.On("GetUser", "123").Return(entities.UserOrganisation{
		Username:          "A",
		OrganisationNames: pq.StringArray{"A"},
	}, nil)
	middleware := middleware.ApiKeyAuthorisation(mockUserRepo)
	gin.SetMode(gin.TestMode)

	responseRecorder := httptest.NewRecorder()
	_, router := gin.CreateTestContext(responseRecorder)
	router.Use(middleware, func(ctx *gin.Context) {
		username := ctx.GetString("username")
		organisationNames := ctx.GetStringSlice("organisations")
		assert.Equal(t, username, "A")
		assert.Equal(t, organisationNames, []string{"A"})
		ctx.JSON(http.StatusOK, gin.H{"username": username, "organisations": organisationNames})
	})

	//given
	authenticatedRequest, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte("{}")))
	authenticatedRequest.Header.Set("Authorization", "Bearer 123")

	//when
	router.ServeHTTP(responseRecorder, authenticatedRequest)

	//then
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	mockUserRepo.AssertExpectations(t)
}

func TestAuthorisationMiddlewareRaisesError(t *testing.T) {
	mockUserRepo := &mocks.UserRepository{}
	middleware := middleware.ApiKeyAuthorisation(mockUserRepo)
	gin.SetMode(gin.TestMode)

	t.Run("invalid authorisation prefix provided", func(t *testing.T) {
		responseRecorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(responseRecorder)
		router.Use(middleware)

		//given
		authenticatedRequest, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte("{}")))
		authenticatedRequest.Header.Set("Authorization", "invalid 123")

		//when
		router.ServeHTTP(responseRecorder, authenticatedRequest)

		//then
		assert.Equal(t, http.StatusForbidden, responseRecorder.Code)
	})

	t.Run("no api key provided", func(t *testing.T) {
		responseRecorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(responseRecorder)
		router.Use(middleware)

		//given
		authenticatedRequest, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte("{}")))

		//when
		router.ServeHTTP(responseRecorder, authenticatedRequest)

		//then
		assert.Equal(t, http.StatusForbidden, responseRecorder.Code)
	})

	t.Run("missing divider", func(t *testing.T) {
		responseRecorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(responseRecorder)
		router.Use(middleware)

		//given
		authenticatedRequest, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte("{}")))
		authenticatedRequest.Header.Set("Authorization", "Bearer123")

		//when
		router.ServeHTTP(responseRecorder, authenticatedRequest)

		//then
		assert.Equal(t, http.StatusForbidden, responseRecorder.Code)
	})
}
