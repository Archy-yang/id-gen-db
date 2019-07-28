package client

import "database/sql"
import _ "github.com/go-sql-driver/mysql"

var mysqlDB *sql.DB

func Init() {
	db, err := sql.Open("mysql", "root:asdfasdf@tcp(127.0.0.1:3306)/id_gen")
	if err != nil {
		panic(err)
	}

	mysqlDB = db
}

func GetMysqlDb() (*sql.DB) {
	return mysqlDB
}
