package mysql

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// Db mysql connection pool
var Db *sql.DB

func init() {
	addr := getEnv("YXI_BACK_MYSQL_ADDR", "127.0.0.1:3306")
	name := getEnv("YXI_BACK_MYSQL_NAME", "yxi")
	user := getEnv("YXI_BACK_MYSQL_USER", "root")
	pass := getEnv("YXI_BACK_MYSQL_PASS", "111")

	var err error
	Db, err = sql.Open("mysql", user+":"+pass+"@tcp("+addr+")/"+name)
	if err != nil {
		log.Fatal(err.Error())
	}

	Db.SetMaxOpenConns(1000)
	Db.SetMaxIdleConns(500)

	err = Db.Ping()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
