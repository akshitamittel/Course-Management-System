package main

import (
	"encoding/json"
	"net/http"
	//"fmt"
)

type queries struct {
	CourseID string `json:"courseID"`
	CourseName string `json:"courseName"`
	Role int `json:"Role"`
}

type context struct {
	Name string
	Email string
	RollNo string
	PL int
	Role string
}

func profileHandler(writer http.ResponseWriter, request *http.Request) {
	var accountID = getAccFromCookie(writer, request, true)
	if accountID == -1 {
		return
	}

	var firstName string
	var lastName string
	var privilegelevel int
	var email string
	var rollno string
	rows , err := mysqldb.Query("select firstname, lastname, privilegelevel, email, rollno from account where accountid=?",accountID)
	if err != nil {
		fatalQueryError("name selection using accountid")
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&firstName, &lastName, &privilegelevel, &email, &rollno)
	}

	// = getJsonResponse()
	if privilegelevel == 2 {
		c := context{Name:(firstName+" "+lastName), Email:email, RollNo:rollno, PL:privilegelevel, Role:"Administrator"}
		templates.ExecuteTemplate(writer, "profile.html", c)
	} else {
		var coursecode string
		var coursename string
		var role int
		rows , err := mysqldb.Query("select coursecode, coursename, role  from course natural join role where accountid=?",accountID)
		if err != nil {
			fatalQueryError("role selection using accountid")
		}
		var obj []queries
		defer rows.Close()
		for rows.Next() {
			rows.Scan(&coursecode, &coursename, &role)
			q := queries{CourseID: coursecode, CourseName: coursename, Role:role}
			obj = append(obj, q)
		}
		// fmt.Println(obj)
		jsonObj, _ := json.Marshal(obj)
		// fmt.Println(string(jsonObj))
		c := context{Name:(firstName+" "+lastName), Email:email, RollNo:rollno, PL:privilegelevel, Role:string(jsonObj)}
		templates.ExecuteTemplate(writer, "profile.html", c)
	}	
}
