package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"fmt"
)

type replyRoleAssn struct {
	Message string
	CourseID int
	MessageID int
}

type previousRole struct {
	Name string `json:"Name"`
	RoleOld int `json:"RoleOld"`
	RoleNew int `json:"RoleNew"`
}

func roleAssnHandler(writer http.ResponseWriter, request *http.Request) {
	var accountID = getAccFromCookie(writer, request, true)
	if accountID == -1 {
		return
	}

	c_id := request.FormValue("courseID")
	fmt.Println(c_id)

	eid := request.FormValue("email")
	fmt.Println(eid)

	rol1 := request.FormValue("role")
	fmt.Println(rol1)

	cid,_ := strconv.Atoi(c_id)
	rol,_ := strconv.Atoi(rol1)
	//Check if Email id exists:
	var d int
	rows, err := mysqldb.Query("select exists (select accountid from account where email=?)",eid)
		if err != nil {
		fatalQueryError("checking for email ID")
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&d)
	}

	//If email id is invalid
	if(d==0){
		fmt.Println("email does not exist")
		c:=replyRoleAssn{Message:"",CourseID:cid,MessageID:-1}
		templates.ExecuteTemplate(writer, "roleRedirect.html", c)
	}else{
		//If email id is valid
		var firstname string
		var lastname string
		var accountid int
		var role int

		rows , err := mysqldb.Query("select accountid, firstname, lastname, role from role natural join account where email=? and courseid=?",eid,cid)
		if err != nil {
			fatalQueryError("checking existing roles in course")
		}
		defer rows.Close()
		if rows.Next() {
			rows.Scan(&accountid,&firstname,&lastname,&role)
			rR:=previousRole{Name:firstname+" "+lastname, RoleOld:role, RoleNew:rol}
			jsonObj, _ := json.Marshal(rR)
			c:=replyRoleAssn{Message:string(jsonObj),CourseID:cid,MessageID:1}

			if(rol != 0){
				var e int
				rows1 , err1 := mysqldb.Query("select accountid from account where email=?",eid)
				if err1 != nil {
					fatalQueryError("retriving accountid from email")
				}
				defer rows1.Close()
				if rows1.Next() {
					rows1.Scan(&e)
				}
				stmt, err1 := mysqldb.Prepare("update role set role=? where courseid=? and accountid=?")
				if err1 != nil {
					fatalQueryError("Could not update role.")
				}

				_, err2 := stmt.Exec(rol, cid, e)
				if err2 != nil {
					fatalQueryError("Could not update role. Problem with executing the update")
				}
			}
			templates.ExecuteTemplate(writer, "roleRedirect.html", c)


		}else{
			c:=replyRoleAssn{Message:"",CourseID:cid,MessageID:0}
			if(rol != 0){
				var e int
				rows1 , err1 := mysqldb.Query("select accountid from account where email=?",eid)
				if err1 != nil {
					fatalQueryError("retriving accountid from email")
				}
				defer rows1.Close()
				if rows1.Next() {
					rows1.Scan(&e)
				}
				mysqldb.Exec("insert into role(accountid,courseid,role) values(?,?,?)",e, cid, rol)
			}
			templates.ExecuteTemplate(writer, "roleRedirect.html", c)
		}
	}
}