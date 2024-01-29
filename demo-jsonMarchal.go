package main

import (
	"encoding/json"
	"fmt"
)

type employee struct {
	ID           int
	EmployeeName string
	Tel          string
	Email        string
}

// ข้อมูลเดิมเป็น object ได้ถูก map ม่เป็นข้อมูลแบบ json
// ข้อควรระวัง key จะต้องเป็นพิมพ์ใหญ่เท่านั้นเช่น ID, EmployeeName, Tel, Email
func main() {
	data,_ := json.Marshal(&employee{101, "Chinnawat", "0966666666","email@email.com"})
	fmt.Println(string(data))
}