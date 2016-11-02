package main

import (
	"net/http"
	"regexp"
	"crypto/sha256"
	"crypto/rand"
	"encoding/hex"
)

func registerHandler(writer http.ResponseWriter, request *http.Request) {
	if getAccFromCookie(writer, request, false) >= 0 {
		//a valid account id is associated with this session id, so go to homepage
		http.Redirect(writer, request, "/home/", http.StatusFound)
		return
	}
	context := struct {
		Emailerr string
	} {}
	if request.URL.Path == "/register/bademail/" {
		context.Emailerr = "Invalid Email address"	
	} else if request.URL.Path == "/register/usedemail/" {
		context.Emailerr = "Email address already in use"
	}
	templates.ExecuteTemplate(writer, "register.html", context)
}

func registerSubmitHandler(writer http.ResponseWriter, request *http.Request) {
	if getAccFromCookie(writer, request, false) >= 0 {
		//a valid account id is associated with this session id, so go to homepage
		http.Redirect(writer, request, "/home/", http.StatusFound)
		return
	}
	firstName := request.FormValue("firstname")
	lastName := request.FormValue("lastname")
	rollNo := request.FormValue("rollno")
	email := request.FormValue("email")
	passwd := request.FormValue("passwd")
	//check that email address is valid
	var validEmail = regexp.MustCompile("^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~])+(\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e])|(\\([\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d)|(([a-zA-Z]|\\d)([a-zA-Z]|\\d|-|\\.|_|~)*([a-zA-Z]|\\d)))\\.)+(([a-zA-Z])|(([a-zA-Z])([a-zA-Z]|\\d|-|\\.|_|~)*([a-zA-Z])))\\.?$")
	m := validEmail.FindStringSubmatch(email)
	if m == nil {
		http.Redirect(writer, request, "/register/bademail/", http.StatusFound)
		return
	}
	rows , err := mysqldb.Query("select accountid from account where email=?",email)
	if err != nil {
		fatalQueryError("account id selection")
	}
	if rows.Next() {
		//an account is already associated with this email, so give an error
		http.Redirect(writer, request, "/register/usedemail/", http.StatusFound)
		rows.Close()
		return
	}
	rows.Close()
	//if no accounts in database, then make this first account a dean account, else make it a student account
	rows , err = mysqldb.Query("select count(accountid) from account")
	defer rows.Close()
	var count int
	if rows.Next() {
		rows.Scan(&count)
	}
	var privilegelevel int
	if count == 0 {
		privilegelevel = 2
	} else {
		privilegelevel = 0
	}
	//encode password
	salt := make([]byte, 16)
	_, err = rand.Read(salt)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	saltstr := hex.EncodeToString(salt)
	passhash := sha256.Sum256( append( []byte(passwd) , salt... ) )
	passhashstr := hex.EncodeToString(passhash[:])
	//add user to database
	mysqldb.Exec("insert into account(firstname,lastname,email,passhash,salt,rollno,privilegelevel) values(?,?,?,?,?,?,?)",
	firstName,lastName,email,passhashstr,saltstr,rollNo,privilegelevel)
	//redirect to login page
	http.Redirect(writer, request, "/login/regcomp/", http.StatusFound)
}
