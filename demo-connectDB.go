package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func query(db*sql.DB){
	var (
		id 			int
		coursName 	string
		price 		float64
		instructor 	string
	)

	for {
		var inputID int
		fmt.Scan(&inputID)

		query := "SELECT id, courseName, price, instructor FROM onlinecourse WHERE id = ? "
		if err := db.QueryRow(query,inputID).Scan(&id, &coursName, &price, &instructor); err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, coursName, price, instructor)
	}
}

// เชื่อมต่อฐานข้อมูล
func main() {
	db,err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/coursedb3")

	if err!=nil {
		fmt.Println("Fail to connect")
	} else {
		fmt.Println("Connect Success")
	}
	fmt.Println(db)
	query(db)
}