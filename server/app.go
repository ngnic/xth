package server

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"xendit-takehome/github/controllers"
	"xendit-takehome/github/middleware"
	"xendit-takehome/github/repositories"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/lib/pq"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type App struct {
	router  *gin.Engine
	appPort string
}

func loadSecretsFromParamStore() {
	region := "ap-southeast-1"
	sess, err := session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{Region: aws.String(region)},
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		panic(err)
	}

	ssmsvc := ssm.New(sess, aws.NewConfig().WithRegion(region))

	withDecryption := true
	params, err := ssmsvc.GetParameters(&ssm.GetParametersInput{
		Names:          []*string{aws.String("DB_URL")},
		WithDecryption: &withDecryption,
	})
	if err != nil {
		panic(err)
	}

	for _, param := range params.Parameters {
		name := aws.StringValue(param.Name)
		value := aws.StringValue(param.Value)
		key := name[strings.LastIndex(name, "/"):]
		os.Setenv(key, value)
	}
}

func migrateSchema(dbType string, db *sql.DB) {
	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		panic(err)
	}

	migrateInstance, err := migrate.NewWithDatabaseInstance(migrationsDir, dbType, driver)
	if err != nil {
		panic(err)
	}

	fmt.Println("Migrating Schema")
	if err := migrateInstance.Up(); err != nil {
		panic(err)
	}
}

func NewApp() App {
	if strings.ToLower(os.Getenv("LOAD_SECRETS_FROM_PARAMSTORE")) == "true" {
		loadSecretsFromParamStore()
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		panic("DB_URL must be specified")
	}
	dbType := dbUrl[:strings.Index(dbUrl, ":")]
	db := sqlx.MustConnect(dbType, dbUrl)

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		panic("APP_PORT must be specified")
	}

	if strings.ToLower(os.Getenv("MIGRATE_DATABASE")) != "false" {
		migrateSchema(dbType, db.DB)
	}

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
