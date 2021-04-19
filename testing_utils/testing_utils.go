package testing_utils

import (
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
)

func GetDBHandle() *sqlx.DB {
	dbUrl := os.Getenv("DB_URL")
	dbType := dbUrl[:strings.Index(dbUrl, ":")]
	return sqlx.MustConnect(dbType, dbUrl)
}

func CleanupTables(db *sqlx.DB) {
	db.MustExec("truncate users cascade")
	db.MustExec("truncate organisations cascade")
}
