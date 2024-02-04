package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB
const coursePath = "courses"
const basePath = "/api"

type Coursed struct {
	CourseID   int     `json:"courseid"`
	CourseName string  `json:"coursename"`
	Price      float64 `json:"price"`
	ImageURL   string  `json:"image_url"`
}

// เชื่อมต่อฐานข้อมูล MySQL
func setupDB() {
	// sql.Open ใช้เพื่อเปิดการเชื่อมต่อกับ MySQL database
	var err error
	DB, err = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/coursedb3")
	if err != nil {
		log.Fatal(err)
	}
	defer DB.Close() // ปิดเมื่อไม่ได้ใช้งาน
	fmt.Println(DB)
	DB.SetConnMaxLifetime(time.Minute * 3)
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(10)
}

// ดึงรายการคอร์สทั้งหมดจากฐานข้อมูล
func getCourseList() ([]Coursed, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// query ข้อมูลทั้งหมดจากฐานข้อมูล
	result, err := DB.QueryContext(ctx, `SELECT
	courseid,
	coursename,
	price,
	image_url
	FROM courseonline`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	defer result.Close()
	courses := make([]Coursed, 0)
	for result.Next() {
		var course Coursed
		result.Scan(
			&course.CourseID,
			&course.CourseName,
			&course.Price,
			&course.ImageURL,
		)
		courses = append(courses, course)
	}
	return courses, nil
}

// แทรกข้อมูลคอร์ส
func insertProduct(course Coursed) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := DB.ExecContext(ctx, `INSERT INTO courseonline
	(
		courseid,
		coursename,
		price,
		image_url
	) VALUE (?, ?, ?, ?)`,
		course.CourseID,
		course.CourseName,
		course.Price,
		course.ImageURL,
	)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	insertID, err := result.LastInsertId()
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}

	return int(insertID), nil
}

// ลบข้อมูล
func removeProduct(courseID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := DB.ExecContext(ctx, `DELETE FROM courseonline where courseid = ?`, courseID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

// ปรับปรุงคอร์ส
func getCourse(courseID int) (*Coursed, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	row := DB.QueryRowContext(ctx, `SELECT
	courseid,
	coursename,
	price,
	image_url
	FROM courseonline WHERE courseid = ?`, courseID)

	course := &Coursed{}
	err := row.Scan(
		&course.CourseID,
		&course.CourseName,
		&course.Price,
		&course.ImageURL,
	)

	// ไม่พบข้อมูล
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return course, nil
}

// จัดการ request ที่เข้ามาที่ endpoint นี้
func handlerCourse(w http.ResponseWriter, r *http.Request) {
	urlPathSegment := strings.Split(r.URL.Path, fmt.Sprintf("%s/", coursePath))

	// [1:] => เอา index ตัวที่ 1 เป็นต้นไป ไม่เอาแค่ 0
	if len(urlPathSegment[1:]) > 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// strconv.Atoi ใช้เพื่อทำการแปลง string เป็น integer.
	// แปลงข้อมูลที่มาจาก URL path segment (โดยเอาตัวสุดท้าย) จาก string เป็น integer (int).
	courseID, err := strconv.Atoi(urlPathSegment[len(urlPathSegment)-1])

	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch r.Method {

	// GET
	case http.MethodGet:
		course, err := getCourse(courseID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if course == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		js, err := json.Marshal(course)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err = w.Write(js)
		if err != nil {
			log.Fatal(err)
		}

	// Delete
	case http.MethodDelete:
		err := removeProduct(courseID)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// จัดการ request ที่เข้ามาที่ endpoint นี้
func handlerCourses(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	// GET
	case http.MethodGet:
		courseList, err := getCourseList()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// แปลง obj เป็น json
		js, err := json.Marshal(courseList)
		if err != nil {
			log.Fatal(err)
		}

		_, err = w.Write(js)
		if err != nil {
			log.Fatal(err)
		}

	// POST
	case http.MethodPost:
		var course Coursed
		err := json.NewDecoder(r.Body).Decode(&course)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		CourseID, err := insertProduct(course)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf(`{"courseid":%id}`, CourseID)))

	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// Middleware สำหรับการจัดการ CORS
func corsMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization, X-CSRF-Token, Accept-Encoding")
		handler.ServeHTTP(w, r)
	})
}

// ตั้งค่า route สำหรับ API
func SetupRoutes(apiBasePath string) {
	courseHandler := http.HandlerFunc(handlerCourse)
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, coursePath), corsMiddleware(courseHandler))

	coursesHandler := http.HandlerFunc(handlerCourses)
	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, coursePath), corsMiddleware(coursesHandler))
}

func main() {
	setupDB()
	SetupRoutes(basePath)
	log.Fatal(http.ListenAndServe(":5000", nil))
}
