package main

import (
	"encoding/json"
	"net/http"
	//"fmt"
	"strconv"
)


type contextStud struct {
	CourseName string
	CourseCode string
	CourseID int
	AllAcc string
}

type queriesStud struct {
	Name string `json:"Name"`
	Email string `json:"Email"`
	Rollno string `json:"Rollno"`
}

func studentHandler(writer http.ResponseWriter, request *http.Request) {
	var accountID = getAccFromCookie(writer, request, true)
	if accountID == -1 {
		return
	}

	//Get the course ID from the URL path
	var courseID = request.URL.Path[len("/students/"):]
	//fmt.Println(courseID)

	var course_id int
	course_id,_ = strconv.Atoi(courseID)

	//Redirect the page to home, if the accessor is not a teacher of the course
	
	var firstname string
	var lastname string
	var email string
	var rollno string
	var obj[] queriesStud
	rows,err := mysqldb.Query("select firstname, lastname, email, rollno from role natural join account where role = 3 and courseid=?",course_id)
	if err != nil {
		fatalQueryError("retrieving student information for course")
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&firstname, &lastname, &email, &rollno)
		q := queriesStud{Name: (firstname+" "+lastname), Email: email, Rollno:rollno}
		obj = append(obj, q)
	}

	//Get course code and course name
	var coursename string
	var coursecode string
	rows , err = mysqldb.Query("select coursename, coursecode from course where courseid=?",course_id)
	if err != nil {
		fatalQueryError("Getting course name from id")
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&coursename, &coursecode)
	}

	jsonObj, _ := json.Marshal(obj)
	//fmt.Println(string(jsonObj))

	c := contextStud{CourseName:coursename, CourseCode:coursecode, CourseID:course_id, AllAcc:string(jsonObj)}
	templates.ExecuteTemplate(writer, "students.html", c)
	
}