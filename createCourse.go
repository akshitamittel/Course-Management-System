package main

import (
	"net/http"
	"html/template"
	"strconv"
)

func createCourseHandler(writer http.ResponseWriter, request *http.Request) {
	var accountID = getAccFromCookie(writer, request, true)
	if accountID == -1 {
		return
	}


	rows,err := mysqldb.Query("select privilegelevel from account where accountid="+strconv.Itoa(accountID))
	
	if err != nil {
		fatalQueryError("privilegelevel selection using accountid")
	}

	defer rows.Close()
	var privilegelevel int
	if rows.Next() {
		rows.Scan(&privilegelevel)
	}
	/* If a student tries to create a page*/
	// if(privilegelevel==0){
	// 	http.Redirect(writer, request, "/home/", http.StatusFound)
	// }



	context := struct {
		Name string
		Email string
		RollNo string
		PL int
		Role string
		CreationComp template.HTML
	} {}

	/*
	NOTE: Not working at the moment. Someone please look into this.
	Successful course creation context
	*/

	if request.URL.Path == "/createCourse/creationcomp/" {
		context.CreationComp = template.HTML("<h2><span style=\"color: green\">Course Creation Successful!</span></h2>")
	}
	templates.ExecuteTemplate(writer, "createCourse.html", context)
}

func createCourseSubmitHandler(writer http.ResponseWriter, request *http.Request) {

	var accountID = getAccFromCookie(writer, request, true)
	if accountID == -1 {
		return
	}

	courseName := request.FormValue("courseName")
	courseCode := request.FormValue("courseCode")
	startDate := request.FormValue("startDate")
	endDate := request.FormValue("endDate")
	maxGraceDays := request.FormValue("maxGraceDays")

	//add course to database
	mysqldb.Exec("insert into course(coursename,coursecode,creatorid,startdate,enddate,maxgracedays) values(?,?,?,?,?,?)",
	courseName,courseCode,accountID,startDate,endDate,maxGraceDays)
	//redirection
	http.Redirect(writer, request, "/createCourse/creationcomp/", http.StatusFound)
}
