package main

import (
	//"encoding/json"
	"net/http"
	// "fmt"
	"strconv"
)


type contextAssnCreate struct {
	CourseName string
	CourseCode string
	CourseID int
}

func assnCreateHandler(writer http.ResponseWriter, request *http.Request) {
	var accountID = getAccFromCookie(writer, request, true)
	if accountID == -1 {
		return
	}

	//Get the course ID from the URL path
	var courseID = request.URL.Path[len("/assnCreate/"):]
	//fmt.Println(courseID)

	var course_id int
	course_id,_ = strconv.Atoi(courseID)

	//Redirect the page to home, if the accessor is not a teacher of the course
	var accountid int
	var teachers int
	teachers = 0
	rows,err := mysqldb.Query("select accountid from role where (role = 1 or role = 2)  and courseid=?",course_id)
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
	templates.ExecuteTemplate(writer, "assnCreate.html", c)
	
}

func assnRegisterHandler(writer http.ResponseWriter, request *http.Request) {
	var accountID = getAccFromCookie(writer, request, true)
	if accountID == -1 {
		return
	}

	c_id := request.FormValue("courseID")
	// fmt.Println(c_id)

	c_type := request.FormValue("createType")
	// fmt.Println(c_type)

	if(c_type == "1"){
		date := request.FormValue("date1")
		time := request.FormValue("time1")
		timedate := date+" "+time+":00"
		date1 := request.FormValue("date2")
		time1 := request.FormValue("time2")
		titlestring := request.FormValue("title")
		marks := request.FormValue("marks")
		var timedate1 string
		if(date1 == ""){
			mysqldb.Exec("insert into assignment(courseid,creatorid,creationtime,duetime,titlestring,submissiontype,maxmarks) values(?,?,NOW(),?,?,\"text\",?)",c_id,accountID,timedate,titlestring,marks)
		}else{
			timedate1 = date1+" "+time1+":00"
			mysqldb.Exec("insert into assignment(courseid,creatorid,creationtime,duetime,maxsubmittime,titlestring,submissiontype,maxmarks) values(?,?,NOW(),?,?,?,\"text\",?)",c_id,accountID,timedate,timedate1,titlestring,marks)
		}
	}

	if(c_type == "2"){
		titlestring := request.FormValue("title")
		mysqldb.Exec("insert into assignment(courseid,creatorid,creationtime,titlestring,submissiontype,maxmarks) values(?,?,NOW(),?,\"none\",0)",c_id,accountID,titlestring)

	}

	http.Redirect(writer, request, "/home/", http.StatusFound)
}