package main

import (
	"net/http"
	"encoding/json"
	"fmt"
)

type queriesHome struct {
	CourseCode string `json:"courseCode"`
	CourseName string `json:"courseName"`
	CourseID int `json:"courseID"`
	Assignments []string `json:"assignments"`
}

type contextHome struct {
	Name string
	Courses string
}

func homeHandler(writer http.ResponseWriter, request *http.Request) {
	var accountID = getAccFromCookie(writer, request, true)
	if accountID == -1 {
		return
	}
	var firstName string
	var lastName string
	var privilegelevel int
	rows , err := mysqldb.Query("select firstname, lastname, privilegelevel from account where accountid=?",accountID)
	if err != nil {
		fatalQueryError("name selection using accountid")
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&firstName,&lastName,&privilegelevel)
	}
	if privilegelevel == 2 {
		//redirect to dean home
		http.Redirect(writer, request, "/deanhome/", http.StatusFound)
	}

	//Get all relevent courseid from roles. Store in a slice
	var courseid int
	var C_ID[] int
	rows , err = mysqldb.Query("select courseid from role where accountid=?",accountID)
	if err != nil {
		fatalQueryError("course selection using accountid")
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&courseid)
		C_ID = append(C_ID, courseid)
	}

	//Using all courseid, get course name, store in slice.
	var coursename string
	var coursecode string
	var C_N[] string
	var C_C[] string
	var obj[] queriesHome
	for _,element := range C_ID{
		rows, err = mysqldb.Query("select coursename, coursecode from course where courseid=?",element)
		if err != nil {
			fatalQueryError("course selection using accountid")
		}
		defer rows.Close()
		if rows.Next() {
			rows.Scan(&coursename, &coursecode)
			C_N = append(C_N, coursename)
			C_C = append(C_C, coursecode)
		}

		var assn[] string
		var titlestring string
		rows, err = mysqldb.Query("select titlestring from assignment where (duetime > NOW() or duetime is null)  and courseid=? order by creationtime desc",element)
		if err != nil {
			fatalQueryError("assignment retreival using courseid in home page")
		}
		defer rows.Close()
		for rows.Next() {
			rows.Scan(&titlestring)
			assn = append(assn, titlestring)
		}

		q := queriesHome{CourseName: coursename, CourseCode: coursecode, CourseID: element, Assignments: assn}
		obj = append(obj, q)

	}
	jsonObj, _ := json.Marshal(obj)
	fmt.Println(string(jsonObj))

	c := contextHome{Name:(firstName+" "+lastName), Courses: string(jsonObj)}

	//[{coursename, coursecode, courseid, {assignments*}}]

	// context := struct {
	// 	Name string
	// } {Name:(firstName+" "+lastName)}
	templates.ExecuteTemplate(writer, "home.html", c)

}

func deanhomeHandler(writer http.ResponseWriter, request *http.Request) {
	var accountID = getAccFromCookie(writer, request, true)
	if accountID == -1 {
		return
	}
	var firstName string
	var lastName string
	var privilegelevel int
	rows , err := mysqldb.Query("select firstname, lastname, privilegelevel from account where accountid=?",accountID)
	if err != nil {
		fatalQueryError("name selection using accountid")
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&firstName,&lastName,&privilegelevel)
	}
	if privilegelevel != 2 {
		//redirect to normal home
		http.Redirect(writer, request, "/home/", http.StatusFound)
	}

	var firstName1 string
	var lastName1 string
	rows , err = mysqldb.Query("select firstname, lastname, privilegelevel from account")
	defer rows.Close()
	//var list_of_names[] string
	var priv_level int
	for rows.Next(){
		rows.Scan(&firstName1, &lastName1, &priv_level)
	}
	context := struct {
		Name string
	} {Name:(firstName+" "+lastName)}
	templates.ExecuteTemplate(writer, "deanhome.html", context)
}
