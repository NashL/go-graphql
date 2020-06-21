package dbConnection

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

var DB *sql.DB

func generateDataSourceName() string{
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	database := os.Getenv("MYSQL_DATABASE")
	sourceName:= user + ":" + password + "@tcp(127.0.0.1:3306)/" + database + "?parseTime=true&sql_mode=ansi"
	return sourceName
}

func NewDatabase() {
	db, err := sql.Open("mysql", generateDataSourceName())
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	DB = db

}
