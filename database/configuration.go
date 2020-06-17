package dbConnection

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)


var DB *sql.DB

func NewDatabase() {
	db, err := sql.Open("mysql", "root:nashwick@tcp(127.0.0.1:3306)/online_store_server?parseTime=true&sql_mode=ansi")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	DB = db
	//defer db.Close()

	// Open doesn't open a connection. Validate DSN data:
	//err = db.Ping()
	//if err != nil {
	//	panic(err.Error()) // proper error handling instead of panic in your app
	//}

	// Use the DB normally, execute the querys etc
	//[...]
}
