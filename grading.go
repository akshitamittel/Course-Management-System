package main

import (
	"net/http"
	"strings"
	"strconv"
	"encoding/json"
)

type submissionData struct {
	accID int
	gradeStatus string
}

type overviewContext struct {
	Title string
	Submissions string
}

type individualContext struct {
	Title string
	AccID int
	LateDays int
	Answer string
	MaxMarks int
}

func gradeHandler(writer http.ResponseWriter, request *http.Request) {
	var accountID = getAccFromCookie(writer, request, true)
	if accountID == -1 {
		return
	}
	parts := strings.Split(request.URL.Path,"/")
	if len(parts)<3 {
		http.Redirect(writer, request, "/home/", http.StatusFound)
		return
	}
	assnID, err := strconv.Atoi(parts[2])
	if err!=nil {
		http.Redirect(writer, request, "/home/", http.StatusFound)
		return
	}
	accID := -1
	if len(parts)==4 {
		accID, err = strconv.Atoi(parts[3])
		if err!=nil {
			http.Redirect(writer, request, "/home/", http.StatusFound)
			return
		}
	}
	courseID := -1
	role	:= -1
	rows , err := mysqldb.Query("select courseid from assignment where assignmentid=?",assnID)
	if err != nil {
		fatalQueryError("course selection from assignment id")
	}
	if rows.Next() {
		rows.Scan(&courseID)
		rows.Close()
	} else {
		fatalQueryError("invalid assignment id")
	}
	rows , err = mysqldb.Query("select role from role where accountid=? and courseid=?",accountID,courseID)
	if err != nil {
		fatalQueryError("role selection from account id and course id")
	}
	if rows.Next() {
		rows.Scan(&role)
		rows.Close()
	} else {
		http.Redirect(writer, request, "/home/", http.StatusFound)
		rows.Close()
		return
	}
	if(role==0) {
		//dont allow grading by students
		http.Redirect(writer, request, "/home/", http.StatusFound)
		return
	}
	if accID==-1 {
		gradeOverviewHandler(writer,request,assnID,courseID)
	} else if len(parts)==4 {
		gradeIndividualHandler(writer,request,assnID,courseID,accID)
	} else {
		http.Redirect(writer, request, "/home/", http.StatusFound)
	}
}

func gradeOverviewHandler(writer http.ResponseWriter, request *http.Request, assnID int, courseID int) {
	var accountID = getAccFromCookie(writer, request, true)
	if accountID == -1 {
		return
	}
	var title string
	var maxMarks int
	rows, err := mysqldb.Query("select titlestring, maxmarks from assignment where assignmentid=?",assnID)
	if err != nil {
		fatalQueryError("title selection from assignment id")
	}
	if rows.Next() {
		rows.Scan(&title,&maxMarks)
		rows.Close()
	} else {
		http.Redirect(writer, request, "/home/", http.StatusFound)
		return
	}
	rows, err = mysqldb.Query("select accountid from role where courseid=? and role=0",courseID)
	if err != nil {
		fatalQueryError("student selection from course id")
	}
	var submissions[] submissionData
	for rows.Next() {
		var accID int
		var sd submissionData
		givenMarks := -1
		rows.Scan(&accID)
		innerRows, err := mysqldb.Query("select givenmarks from grade where assignmentid=? and submittorid=?",assnID,accID)
		if err != nil {
			fatalQueryError("grade selection from assignment and submittor id")
		}
		if innerRows.Next() {
			innerRows.Scan(&givenMarks)
			sd = submissionData{accID:accID, gradeStatus:(strconv.Itoa(givenMarks)+"/"+strconv.Itoa(maxMarks))}
		} else {
			sd = submissionData{accID:accID, gradeStatus:"Not Graded"}
		}
		innerRows.Close()
		submissions = append(submissions,sd)
	}
	jsonObj, _ := json.Marshal(submissions)
	ctxt := overviewContext{Title:title, Submissions:string(jsonObj)}
	templates.ExecuteTemplate(writer, "gradeOverview.html", ctxt)
}

func gradeIndividualHandler(writer http.ResponseWriter, request *http.Request, assnID int, courseID int, accID int) {
	var accountID = getAccFromCookie(writer, request, true)
	if accountID == -1 {
		return
	}
	var title string
	var maxMarks int
	var dueTime string
	rows, err := mysqldb.Query("select titlestring, maxmarks, duetime from assignment where assignmentid=?",assnID)
	if err != nil {
		fatalQueryError("title selection from assignment id")
	}
	if rows.Next() {
		rows.Scan(&title,&maxMarks,&dueTime)
		rows.Close()
	} else {
		http.Redirect(writer, request, "/home/", http.StatusFound)
		return
	}
	var answer string
	var submitTime string
	rows, err = mysqldb.Query("select answerstring, submittime from submstring where submittorid=? and assignmentid=?",accID,assnID)
	if err != nil {
		fatalQueryError("getting submitted answer")
	}
	if rows.Next() {
		rows.Scan(&answer,&submitTime)
		rows.Close()
	} else {
		http.Redirect(writer, request, "/home/", http.StatusFound)
		return
	}
	ctxt := individualContext{Title:title,AccID:accID,LateDays:0,Answer:answer,MaxMarks:maxMarks}
	templates.ExecuteTemplate(writer, "grade.html", ctxt)
}

func gradeSubmitHandler(writer http.ResponseWriter, request *http.Request) {
	var accountID = getAccFromCookie(writer, request, true)
	if accountID == -1 {
		return
	}
	assnID, err := strconv.Atoi(request.FormValue("assnID"))
	if err != nil {
		http.Redirect(writer, request, "/home/", http.StatusFound)
		return
	}
	accID, err := strconv.Atoi(request.FormValue("accID"))
	if err != nil {
		http.Redirect(writer, request, "/home/", http.StatusFound)
		return
	}
	givenMarks, err := strconv.Atoi(request.FormValue("given"))
	if err != nil {
		http.Redirect(writer, request, "/home/", http.StatusFound)
		return
	}
	mysqldb.Exec("insert into grade(assignmentid,submittorid,givenmarks,gracedaysused) values (?,?,?,0)",assnID,accID,givenMarks)
	http.Redirect(writer, request, "/grade/"+strconv.Itoa(assnID), http.StatusFound)
}