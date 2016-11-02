package main

import (
	"net/http"
	"fmt"
	"time"
	"strconv"
	"encoding/json"
)


type contextAssignment struct {
	AssnName string 
	AssnID int
	Instructors string
	CourseCode string
	Creationtime string
	Duetime string
	MaxMarks string
	Role int
	Total int 
	Done string
	Maxsubmittime string
	Comments string

}

type qfiles struct{
	Fileid int `json:"fid"`
	Filename string `json:"fname"`
}

type comments struct{
	Commentstring string `json:"cstring"`
	Commentorid int `json:"cid"`
	Commenttime string `json:"ctime"`
}



func assignmentHandler(writer http.ResponseWriter, request *http.Request) {
	var accountID = getAccFromCookie(writer, request, true)
	if accountID == -1 {
		return
	}

	//GET ASSN_ID FROM URL:
	var assnID = request.URL.Path[len("/assn/"):]
	fmt.Println(assnID)
	var assn_id int
	assn_id,_ = strconv.Atoi(assnID)

	//GET ASSIGNMENT INFORMATION:
	var assnname string
	var course_id string
	rows, err := mysqldb.Query("select titlestring, courseid from assignment where assignmentid = ?",assn_id)
	if err != nil {
		fatalQueryError("retrieving AssignmentInfo")
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&assnname, &course_id)
	}

	//GET COURSEID:
	var coursecode string
	rows, err = mysqldb.Query("select coursecode from course where courseid=?",course_id)
	if err != nil {
		fatalQueryError("retrieving CourseID in Assignment page")
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&coursecode)
	}



	//GET INSTRUCTORS AND TA'S OF THE COURSE:
	var firstname string
	var lastname string
	var name string
	var accountid int
	var teachers[] string
	var tas[] string
	var role int
	rows, err = mysqldb.Query("select accountid, firstname, lastname, role from account natural join role  where courseid=?",course_id)
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

	instr := instructors{Teachers:teachers}
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

	//GET OTHER DETAILS:
	var creationtime string
	var duetime string
	var maxsubmittime string
	var maxmarks string	
	var done string
	var d int
	var count int



	rows, err = mysqldb.Query("select creationtime, duetime, maxsubmittime, maxmarks from assignment where assignmentid=?",assn_id)
	if err != nil {
		fatalQueryError("retrieving Assignment Info")
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&creationtime, &duetime, &maxsubmittime, &maxmarks)

		if(role==1 || role == 2){
			rows1, err1 := mysqldb.Query("select count(distinct p.submittorid) from (select submittorid, assignmentid from submstring union select submittorid, assignmentid from subfile)p where p.assignmentid=?",assn_id)
			if err1 != nil {
				fmt.Println(err1)
				fatalQueryError("Assignment completion status in Assignment page")
			}
			defer rows1.Close()
			if rows1.Next() {
				rows1.Scan(&count)
			}
			done = "N/A"
		}else{
			if(duetime != ""){
				rows1, err1 := mysqldb.Query("select exists(select p.* from (select submittorid, assignmentid from submstring union select submittorid, assignmentid from subfile)p where p.submittorid=? and p.assignmentid=?)", accountID, assn_id)
				if err1 != nil {
					fmt.Println(err1)
					fatalQueryError("Assignment completion status in Assignment page")
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
		

	}

	//GET COMMENTS:
	var commentstring string
	var commentorid int
	var commenttime string
	
	var comments_array[] comments
	rows, err = mysqldb.Query("select commentstring, commentorid, commenttime from comment where assignmentid = ? order by commenttime desc",assn_id)
	if err != nil {
		fatalQueryError("retrieving comments in assignment page ")
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&commentstring, &commentorid, &commenttime)
		
		c := comments{Commentstring:commentstring, Commentorid:commentorid, Commenttime:commenttime}
		comments_array = append(comments_array, c)
	}

	jsonObjComment,_ := json.Marshal(comments_array)

	fmt.Println(assn_id)
	fmt.Println(assnname)
	fmt.Println(coursecode)
	fmt.Println(creationtime)
	fmt.Println(duetime)
	fmt.Println(maxmarks)
	fmt.Println(role)
	fmt.Println(count)
	fmt.Println(done)

	//SEND INFORMATION TO THE PAGE:
	a := contextAssignment{AssnName:assnname, AssnID:assn_id, Instructors:string(jsonObjInstr), CourseCode:coursecode, Creationtime:creationtime, Duetime:duetime, MaxMarks:maxmarks, Role:role, Total:count, Done:done, Maxsubmittime:maxsubmittime, Comments:string(jsonObjComment)}
	templates.ExecuteTemplate(writer, "assn.html", a)
	
}


func commentSubmitHandler(writer http.ResponseWriter, request *http.Request) {
	var accountID = getAccFromCookie(writer, request, true)
	if accountID == -1 {
		return
	}

	//GET ASSN_ID FROM URL:
	var assnID = request.FormValue("asn_id")
	//fmt.Println(assnID)
	var assn_id int
	assn_id,_ = strconv.Atoi(assnID)

	//assn_id =1



	var comment string

	comment = request.FormValue("comment")

	if(comment!=""){
	//add answer to database
		current_time := time.Now().Local()

		fmt.Println(comment)
		fmt.Println(accountID)
		fmt.Println(assn_id)
		fmt.Println(current_time)

		mysqldb.Exec("insert into comment(commentstring,commentorid,assignmentid,commenttime) values (?,?,?,?)",
		comment,accountID,assn_id,current_time)

	} 



	var answerstring string

	answerstring = request.FormValue("answer")

	//assn_id =1

	fmt.Println("here")
	
	if(answerstring!=""){
		current_time := time.Now().Local()
		fmt.Println(answerstring)
		fmt.Println(accountID)
		fmt.Println(assn_id)
		fmt.Println(current_time)

		//add answer to database
		mysqldb.Exec("insert into submstring(answerstring,submittorid,assignmentid,submittime) values (?,?,?,?)",
		answerstring,accountID,assn_id,current_time)
		fmt.Println("Wait!!!!!!!!!!!!!!")

	}





   //redirection
	http.Redirect(writer, request, "/assn/"+assnID, http.StatusFound)
}
