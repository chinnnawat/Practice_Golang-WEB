package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
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

// Func Find ID
func findID(ID int)(*Course, int){
	// i = ตำแหน่ง index
	// range จะทำให้เข้าถึงทุก element ใน CourseList
	for i, course := range CourseList{
		if course.CourseId == ID {
			return &course, i
		}
	}
	return nil,0
}

func courseHandler(w http.ResponseWriter, r *http.Request) {

	// ทำการ split request ที่ส่งมา โดยใช้ course/ ในการ splite
	// ตัวอย่าง URL path: "/course/123"
	// เมื่อใช้ strings.Split และ "course/" เป็น delimiter, เราจะได้ slice ของ strings ที่มีสองส่วน: ["", "123"]
	// เรานำตัวสุดท้ายของ slice นี้ (index -1) ที่เป็น "123" มาแปลงเป็นตัวเลขโดยใช้ strconv.Atoi และนำมาใช้เป็น ID ของคอร์สที่เราต้องการ.
	urlPathSegment := strings.Split(r.URL.Path, "course/")

	// ทำการแปลงตัว string
	ID, err := strconv.Atoi(urlPathSegment[len(urlPathSegment)-1])
	if err != nil{
		log.Print(err)
		// show status
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// findID(ID) ส่งค่ามา 2 ตัว คือ 1.&course(Course) 2.i(int) และตัวที่รับคือ course, listItemIndex ตามลำดับ
	course, listItemIndex := findID(ID)
	if course == nil {
		http.Error(w, fmt.Sprintf("no course withd id", ID),http.StatusNotFound)
		return
	}

	switch r.Method {
	// Method Get ใช้ดูข้อมูลเฉพาะ id นั้นๆ
	case http.MethodGet:
		courseJSON, err := json.Marshal(course)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(courseJSON)
	case http.MethodPut:
		var updateCourse Course
		byteBody, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(byteBody, &updateCourse)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if updateCourse.CourseId != ID{
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		course = &updateCourse
		
		//  ทำหน้าที่แทนที่ข้อมูลของคอร์สที่อยู่ใน CourseList ที่ตำแหน่ง listItemIndex 
		// ด้วยข้อมูลใหม่ที่ได้จากการอัปเดต (ที่เก็บไว้ในตัวแปร course).
		CourseList[listItemIndex] = *course
		w.WriteHeader(http.StatusOK)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// ใช้จัดการ req, res ที่เกิดขึ้น
func coursesHandler(w http.ResponseWriter, r *http.Request) {
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
	http.HandleFunc("/course/", courseHandler)
	http.HandleFunc("/course", coursesHandler)
	http.ListenAndServe(":5000", nil)
}