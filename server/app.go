package server

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"xendit-takehome/github/controllers"
	"xendit-takehome/github/middleware"
	"xendit-takehome/github/repositories"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/file"
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

	environment := os.Getenv("ENVIRONMENT")
	withDecryption := true
	params, err := ssmsvc.GetParameters(&ssm.GetParametersInput{
		Names:          []*string{aws.String(fmt.Sprintf("/%s/app/DB_URL", environment))},
		WithDecryption: &withDecryption,
	})
	if err != nil {
		panic(err)
	}

	for _, param := range params.Parameters {
		name := aws.StringValue(param.Name)
		value := aws.StringValue(param.Value)
		key := name[strings.LastIndex(name, "/")+1:]
		os.Setenv(key, value)
	}
}

func migrateSchema(dbType string, db *sql.DB) {
	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		panic(err)
	}

	fsrc, err := (&file.File{}).Open(migrationsDir)
	if err != nil {
		panic(err)
	}

	migrateInstance, err := migrate.NewWithInstance("file", fsrc, dbType, driver)
	if err != nil {
		panic(err)
	}

	if err := migrateInstance.Up(); err != nil && err != migrate.ErrNoChange {
		panic(err)
	}
}

func NewApp() App {
	if strings.ToLower(os.Getenv("LOAD_SECRETS_FROM_PARAMSTORE")) == "true" {
		log.Println("Loading params from paramstore")
		loadSecretsFromParamStore()
		log.Println("Params loaded successfully")
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
		log.Println("Performing migrations")
		migrateSchema(dbType, db.DB)
		log.Println("Migrations were performed successfully")
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
