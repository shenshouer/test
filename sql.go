package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)
	mysql_url := "root:123456@tcp(127.0.0.1:3306)/mediadist?allowAllFiles=true&charset=utf8&interpolateParams=true&parseTime=true"
	db, err := sql.Open("mysql", mysql_url)
	if err != nil {
		log.Fatal(err)
	}

	row := db.QueryRow("SELECT * FROM test1 WHERE id=3")
	var id int
	var name string

	err = row.Scan(&id, &name)

	if err != nil && sql.ErrNoRows != err {
		log.Fatal(err)
	}

	fmt.Println(id, name)
}
