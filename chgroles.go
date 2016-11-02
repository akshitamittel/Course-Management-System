package main

import (
	//"encoding/json"
	"net/http"
	//"fmt"
	"strconv"
)


type contextChg struct {
	CourseName string
	CourseCode string
	CourseID int
}

func courseRolesHandler(writer http.ResponseWriter, request *http.Request) {
	var accountID = getAccFromCookie(writer, request, true)
	if accountID == -1 {
		return
	}

	//Get the course ID from the URL path
	var courseID = request.URL.Path[len("/chgroles/"):]
	//fmt.Println(courseID)

	var course_id int
	course_id,_ = strconv.Atoi(courseID)

	//Redirect the page to home, if the accessor is not a teacher of the course
	var accountid int
	var teachers int
	teachers = 0
	rows,err := mysqldb.Query("select accountid from role where role = 1 and courseid=?",course_id)
	if err != nil {
		fatalQueryError("name selection using accountid")
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&accountid)
		if(accountid == accountID){
			teachers = 1
		}
	}

	//fmt.Println(teachers)
	if(teachers == 0){
		http.Redirect(writer, request, "/home/", http.StatusFound)

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


	c := contextChg{CourseName:coursename, CourseCode:coursecode, CourseID:course_id}
	templates.ExecuteTemplate(writer, "chgroles.html", c)
	
}