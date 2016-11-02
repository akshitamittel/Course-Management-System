package main

import (
	"net/http"
	"strings"
	"strconv"
	"os/exec"
)

func downloadHandler(writer http.ResponseWriter, request *http.Request) {
	var accountID = getAccFromCookie(writer, request, true)
	if accountID == -1 {
		return
	}
	fileid, err := strconv.Atoi(strings.TrimPrefix(request.URL.Path,"/download/"))
	rows , err := mysqldb.Query("select realfname, storedfname from file where fileid=?",fileid)
	if err != nil {
		fatalQueryError("file selection for download")
	}
	var realfname string
	var storedfname string
	if rows.Next() {
		rows.Scan(&realfname,&storedfname)
		rows.Close()
	} else {
		//file not found
		http.Redirect(writer, request, "/home/", http.StatusFound)
		rows.Close()
		return
	}
	var parts = strings.Split(strings.TrimPrefix(storedfname,realfname),"_")
	assnID, err := strconv.Atoi(parts[1])
	var accID = -1
	var courseID = -1
	var role = -1
	if(parts[2]!="q") {
		accID, err = strconv.Atoi(parts[2])
	}
	rows , err = mysqldb.Query("select courseid from assignment where assignmentid=?",assnID)
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
		//don't allow random file download
		http.Redirect(writer, request, "/home/", http.StatusFound)
		rows.Close()
		return
	}
	if(accID==-1 || accID==accountID || role>0) {
		//if question file, or if own submission, or if TA or instructor, then allow download
		var cmd *exec.Cmd
		cmd = exec.Command("/bin/cp","./files/"+storedfname,"./"+realfname)
		err = cmd.Run() //should work, but error can occur here
		http.ServeFile(writer,request,"./"+realfname)
		cmd = exec.Command("/bin/rm","./"+realfname)
		err = cmd.Run() //should work, but error can occur here
	} else {
		//don't allow random file download
		http.Redirect(writer, request, "/home/", http.StatusFound)
	}
}
