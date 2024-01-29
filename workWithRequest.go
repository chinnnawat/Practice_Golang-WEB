package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Course struct {
	CourseId   int    `json:"id"`
	CourseName string `json:"name"`
	Price      int    `json:"price"`
	Instructor string `json:"instructor"`
}

// [] = slice
// คือการประกาศตัวแปร CourseList ซึ่งเป็น slice ของ Course
var CourseList []Course

func init() {
	CourseJSON := `[
		{
			"id": 1,
			"name": "Introduction to Programming",
			"price": 500,
			"instructor": "John Doe"
		},
		{
			"id": 2,
			"name": "Web Development Basics",
			"price": 800,
			"instructor": "Jane Smith"
		},
		{
			"id": 3,
			"name": "Data Structures and Algorithms",
			"price": 1000,
			"instructor": "Bob Johnson"
		},
		{
			"id": 4,
			"name": "Mobile App Development",
			"price": 1200,
			"instructor": "Alice Williams"
		},
		{
			"id": 5,
			"name": "Advanced Machine Learning",
			"price": 1500,
			"instructor": "Charlie Brown"
		}
	]
	`

	// แปลง json เป็น object
	err := json.Unmarshal([]byte(CourseJSON), &CourseList)
	if err != nil {
		log.Fatal(err)
	}
}

// generate next ID
func getNextID()int{
	highestID := 1
	// for _, course := range CourseList วนลูปจนกว่าจะเช็ค object ทุกตัวใน CourseList
	// _ ใช้เพื่อบอกว่าไม่มีการใช้ตัวแปรอะไรที่เกี่ยวข้องกับ for แค่ต้องการใช้ลูป for เฉยๆ 
	// ถ้าคุณไม่ใช้ _ และกำหนดตัวแปร, แต่ไม่ได้ใช้ตัวแปรนั้นภายใน loop, Go 
	// จะแจ้งเตือนว่า "declared and not used" หมายความว่าคุณประกาศตัวแปรไว้แต่ไม่ได้ใช้.
	for _, course := range CourseList {
		if  highestID < course.CourseId{
			highestID = course.CourseId
		}
	}
	return highestID+1
}

// ใช้จัดการ req, res ที่เกิดขึ้น
func courseHandler(w http.ResponseWriter, r *http.Request) {
	// แปลง obj => json
	courseJSON, err := json.Marshal(CourseList)
	switch r.Method{
	// case Get
	case http.MethodGet:
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(courseJSON)
	
	// case post
	case http.MethodPost:
		var newCourse Course
		// คำสั่งอ่านทั้งหมด io
		BodyByte, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// ยัด BodyByte ใน newCourse โดย pointer ไปยัง newCourse
		err = json.Unmarshal(BodyByte, &newCourse)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		newCourse.CourseId = getNextID()
		CourseList = append(CourseList, newCourse)
		w.WriteHeader(http.StatusCreated)
		return
	}
}

func main() {
	http.HandleFunc("/course", courseHandler)
	http.ListenAndServe(":5000", nil)
}