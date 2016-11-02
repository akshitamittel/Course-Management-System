package main

import (
	//"encoding/json"
	"net/http"
	//"fmt"
	"strconv"
)

type accountinfo struct {
	Name string
	PrivilegeLevel int
	Email string
	Rollno string
	Response int
	AccEmail string
}

func privilegeHandler(writer http.ResponseWriter, request *http.Request) {
	var accountID = getAccFromCookie(writer, request, true)
	if accountID == -1 {
		return
	}

	var firstName string
	var lastName string
	var emailid string
	var accountEmail string
	var email string
	var rollno string
	var privilegelevel int
	rows , err := mysqldb.Query("select email, privilegelevel from account where accountid=?",accountID)
	if err != nil {
		fatalQueryError("name selection using accountid")
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&accountEmail,&privilegelevel)
	}
	if privilegelevel != 2 {
		//redirect to normal home
		http.Redirect(writer, request, "/home/", http.StatusFound)
	}
	
	emailid = request.FormValue("email")
	rows , err = mysqldb.Query("select firstname, lastname, privilegelevel, email, rollno from account where email=?",emailid)
	if err != nil {
		fatalQueryError("name selection using accountid")
	}

	if rows.Next() {
		rows.Scan(&firstName,&lastName,&privilegelevel, &email, &rollno)
		//fmt.Println(firstName)
		c := accountinfo{Name:(firstName+" "+lastName), PrivilegeLevel: privilegelevel, Email: email, Rollno: rollno, Response: 1,AccEmail:accountEmail}
		templates.ExecuteTemplate(writer, "privileges.html", c)

	}else{
		c := accountinfo{Name: "None", PrivilegeLevel: -1, Email: "email", Rollno: "rollno", Response: 0, AccEmail: "sup"}
		templates.ExecuteTemplate(writer, "privileges.html", c)	
	}
	

}

func privilegeEscHandler(writer http.ResponseWriter, request *http.Request){
	var accountID = getAccFromCookie(writer, request, true)
	if accountID == -1 {
		return
	}

	chosenLevel, err := strconv.Atoi(request.FormValue("Level"))
	if err != nil {
		fatalQueryError("Problem with level")
	}
	EmailID := request.FormValue("emailID")
	stmt, err1 := mysqldb.Prepare("update account set privilegelevel=? where email=?")
	if err1 != nil {
		fatalQueryError("Could not update privilege level. Problem with emailID or privLevel")
	}

	_, err2 := stmt.Exec(chosenLevel, EmailID)
	if err2 != nil {
		fatalQueryError("Could not update privilege level. Problem with executing the update")
	}

	http.Redirect(writer, request, "/home/", http.StatusFound)

	
}