package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// create ตาราง
func creatingTable(db*sql.DB){
	// `คำสั่งที่ใช้ create table ใน mysql`
	// id INT AUTO_INCREMENT = การสร้าง id (int) โดยอัตโนมัติ
	query := `CREATE TABLE users(
		id INT AUTO_INCREMENT,
		username TEXT NOT NULL,
		password TEXT NOT NULL,
		created_at DATETIME,
		PRIMARY KEY (id)
	);`

	// Exec = execute SQL query
	if _, err := db.Exec(query); err != nil {
		log.Fatal(err)
	}
}

// อ่านข้อมูลท่ได้จากฐานข้อมูล
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

		// "(คำสั่งที่ใช้ใน mysql)"
		query := "SELECT id, courseName, price, instructor FROM onlinecourse WHERE id = ? "

		if err := db.QueryRow(query,inputID).Scan(&id, &coursName, &price, &instructor); err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, coursName, price, instructor)
	}
}

// insert user เข้าสู่ตารางใน mysql
func insert(db*sql.DB){
	var (
		username string
		password string	
	)
	fmt.Scan(&username)
	fmt.Scan(&password)
	createAt := time.Now()
	result,err := db.Exec(`INSERT INTO users (username, password, created_at) VALUES (?, ?, ?);`, username, password, createAt);
	if err != nil {
		log.Fatal(err)
	}
	id, err := result.LastInsertId()
	fmt.Println(id)
}

// delete User
func deleteUserId(db*sql.DB){
	var deleteid int
	fmt.Scan(&deleteid)
	_, err := db.Exec(`DELETE FROM users WHERE id = ?`,deleteid)
	if err != nil {
		log.Fatal(err)
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

	// createTable
	// creatingTable(db)

	// insert User
	// insert(db)

	// deleteUserId
	deleteUserId(db)
}