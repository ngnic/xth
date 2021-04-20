package server

import (
	"fmt"
	"os"
	"strings"
	"xendit-takehome/github/controllers"
	"xendit-takehome/github/middleware"
	"xendit-takehome/github/repositories"

	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type App struct {
	router  *gin.Engine
	appPort string
}

func NewApp() App {
	dbUrl := os.Getenv("DB_URL")
	dbType := dbUrl[:strings.Index(dbUrl, ":")]
	db := sqlx.MustConnect(dbType, dbUrl)

	appPort := os.Getenv("APP_PORT")

	router := gin.Default()
	orgRepo := repositories.NewOrganisationDBRepository(db)
	userRepository := repositories.NewUserDBRepository(db)
	authMiddleware := middleware.ApiKeyAuthorisation(userRepository)
	controllers.SetupRoutes(router, orgRepo, authMiddleware)

	return App{
		router:  router,
		appPort: appPort,
	}
}

func (app *App) Run() {
	app.router.Run(fmt.Sprintf(":%s", app.appPort))
}
