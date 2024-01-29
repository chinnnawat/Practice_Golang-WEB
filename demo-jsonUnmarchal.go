package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type employee struct {
	ID           int
	EmployeeName string
	Tel          string
	Email        string
}

// แปลงจาก json มาเป็น object
func main() {
	e := employee{}
	err := json.Unmarshal([]byte(`{"ID":101, "EmployeeName":"Chinnawat", "Tel":"0966666666", "Email":"email@email.com"}`),&e)
	if err != nil {
		// การใช้ log.Fatal(err) จะแสดงข้อความที่ถูกพิมพ์ด้วย log.Print พร้อมกับ
		// การเรียก os.Exit(1) เพื่อจบโปรแกรมทันที. 
		// ในทางปฏิบัติ, คุณจะใช้ log.Fatal เมื่อคุณต้องการจบโปรแกรมเมื่อเจอข้อผิดพลาดและไม่สามารถดำเนินการต่อได้
		log.Fatal(err)
	}
	fmt.Println(e.EmployeeName)
}