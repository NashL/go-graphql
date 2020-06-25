package dbConnection

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

var DB *sql.DB

func localDataSourceName() string{
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	database := os.Getenv("MYSQL_DATABASE")
	sourceName:= user + ":" + password + "@tcp(127.0.0.1:3306)/" + database + "?parseTime=true&sql_mode=ansi"
	return sourceName
}

func cloudDataSourceName() string {
	var (
		dbUser                 = mustGetenv("DB_USER")
		dbPwd                  = mustGetenv("DB_PASS")
		instanceConnectionName = mustGetenv("INSTANCE_CONNECTION_NAME")
		dbName                 = mustGetenv("DB_NAME")
	)

	var dbURI string = dbUser + ":" + dbPwd + "@unix(/cloudsql/" + instanceConnectionName + ")/" + dbName
	return dbURI
}


func NewDatabase() {

	localEnvironment := os.Getenv("LOCAL_ENVIRONMENT")
	var dbURI string
	if localEnvironment == "false" {
		dbURI = cloudDataSourceName()
	} else {
		dbURI = localDataSourceName()
	}

	db, err := sql.Open("mysql", dbURI)
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	DB = db

}

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Panicf("%s environment variable not set.", k)
	}
	return v
}

