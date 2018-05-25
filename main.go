package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "root:111@tcp(127.0.0.1:3306)/yxi")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalln(err)
	}

	rows, err := db.Query("select * from test1")
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()

	for rows.Next() {
		var tt te
		rows.Scan(&tt.id, &tt.name)
		fmt.Println(tt.id, tt.name)
	}

}

type te struct {
	id   int
	name string
}
