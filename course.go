package main

import (
	"net/http"
	"fmt"
	"strconv"
	"encoding/json"
)


type contextCourse struct {
	CourseName string
	Instructors string
	CourseCode string
	CourseID int
	Role int
	Assignments string
}

type instructors struct {
	Teachers[] string `json:"teachers"`
	Tas[] string `json:"tas"`
}

type assignments struct{
	Assignmentid int `json:"id"`
	Creationtime string `json:"creation"`
	Done string `json:"done"`
	Duetime string `json:"due"`
	Maxsubmittime string `json:"maxDue"`
	Titlestring string `json:"title"`
	Total int `json:"total"`
}

func courseHandler(writer http.ResponseWriter, request *http.Request) {
	var accountID = getAccFromCookie(writer, request, true)
	if accountID == -1 {
		return
	}

	//GET COURSE_ID FROM URL:
	var courseID = request.URL.Path[len("/course/"):]
	//fmt.Println(courseID)
	var course_id int
	course_id,_ = strconv.Atoi(courseID)

	//GET COURSE INFORMATION:
	var coursename string
	var coursecode string
	rows, err := mysqldb.Query("select coursename, coursecode from course where courseid=?",course_id)
	if err != nil {
		fatalQueryError("retrieving courseInfo")
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&coursename, &coursecode)
	}

	//GET INSTRUCTORS AND TA'S OF THE COURSE:
	var firstname string
	var lastname string
	var name string
	var accountid int
	var teachers[] string
	var tas[] string
	var role int
	rows, err = mysqldb.Query("select accountid, firstname, lastname, role from account natural join role where courseid=?",course_id)
	if err != nil {
		fatalQueryError("retrieving instructor info in course page")
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&accountid, &firstname, &lastname, &role)
		name = firstname+" "+lastname
		if (role == 1){
			teachers = append(teachers, name)
		} else if (role == 2){
			tas = append(tas, name)
		}
	}

	instr := instructors{Teachers:teachers,Tas:tas}
	jsonObjInstr, _ := json.Marshal(instr)

	fmt.Println(string(jsonObjInstr))

	//GET THE ROLE OF THE USER ACCESSING THE PAGE:
	rows , err = mysqldb.Query("select role from role where accountid=? and courseid=?",accountID,course_id)
	if err != nil {
		fatalQueryError("retrieving role")
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&role)
	}

	//GET ASSIGNMENTS/ANNOUNCEMENTS:
	var assignmentid int
	var creationtime string
	var duetime string
	var maxsubmittime string
	var titlestring string
	var done string
	var d int
	var count int
	var assns[] assignments
	rows, err = mysqldb.Query("select titlestring, assignmentid, creationtime, duetime, maxsubmittime from assignment where courseid=? order by creationtime desc",course_id)
	if err != nil {
		fatalQueryError("retrieving assignments in course page ")
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&titlestring, &assignmentid, &creationtime, &duetime, &maxsubmittime)
		fmt.Println("Titlestring:")
		fmt.Println(titlestring)
		if(role==1 || role == 2){
			rows1, err1 := mysqldb.Query("select count(distinct p.submittorid) from (select submittorid, assignmentid from submstring union select submittorid, assignmentid from subfile)p where p.assignmentid=?",assignmentid)
			if err1 != nil {
				fmt.Println(err1)
				fatalQueryError("Assignment completion status in course page")
			}
			defer rows1.Close()
			if rows1.Next() {
				rows1.Scan(&count)
			}
			done = "N/A"
		}else{
			if(duetime != ""){
				rows1, err1 := mysqldb.Query("select exists(select p.* from (select submittorid, assignmentid from submstring union select submittorid, assignmentid from subfile)p where p.submittorid=? and p.assignmentid=?)", accountID, assignmentid)
				if err1 != nil {
					fmt.Println(err1)
					fatalQueryError("Assignment completion status in course page")
				}
				defer rows1.Close()
				if rows1.Next() {
					rows1.Scan(&d)
					if (d == 0) {
						done = "Not Done"
					} else {
						done = "Done"
					}
				}

			} else {
				done = "none"
			}
			count = -1
		}
		a := assignments{Assignmentid:assignmentid, Creationtime:creationtime, Duetime:duetime, Done:done, Maxsubmittime:maxsubmittime, Titlestring:titlestring, Total:count}
		assns = append(assns, a)
	}

	jsonObjAssn,_ := json.Marshal(assns)

	fmt.Println(string(jsonObjAssn))

	

	//GET THE NUMBER OF SUBMISSIONS DONE BY STUDENTS FOR TA'S AND TEACHERS
	rows, err = mysqldb.Query("")

	//SEND INFORMATION TO THE PAGE:
	c := contextCourse{CourseCode:coursecode, Instructors:string(jsonObjInstr), CourseName:coursename, CourseID: course_id, Role:role, Assignments:string(jsonObjAssn)}
	templates.ExecuteTemplate(writer, "course.html", c)
	
}
