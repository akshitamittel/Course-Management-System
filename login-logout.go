package main

import (
	"net/http"
	"html/template"
	"crypto/sha256"
	"crypto/rand"
	"encoding/hex"
)

//handler for the login page
func loginHandler(writer http.ResponseWriter, request *http.Request) {
	if getAccFromCookie(writer, request, false) >= 0 {
		//a valid account id is associated with this session id, so go to homepage
		http.Redirect(writer, request, "/home/", http.StatusFound)
		return
	}
	context := struct {
		Emailerr string
		Passwderr string
		RegComp template.HTML //cannot use string since HTML tags are automatically escaped, we do not want that.
	} {}
	if request.URL.Path == "/login/emailerr/" {
		context.Emailerr = "Invalid or Unregistered Email"	
	} else if request.URL.Path == "/login/passwderr/" {
		context.Passwderr = "Incorrect Password"
	} else if request.URL.Path == "/login/regcomp/" {
		context.RegComp = template.HTML("<h2><span style=\"color: green\">Registration Complete! Login to your new account!</span></h2>")
	}
	templates.ExecuteTemplate(writer, "login.html", context)
}

//handler for login form submissions
func loginSubmitHandler(writer http.ResponseWriter, request *http.Request) {
	if getAccFromCookie(writer, request, false) >= 0 {
		//a valid account id is associated with this session id, so go to homepage
		http.Redirect(writer, request, "/home/", http.StatusFound)
		return
	}
	email := request.FormValue("email")
	passwd := request.FormValue("passwd")
	rows , err := mysqldb.Query("select accountid, passhash, salt from account where email=?",email)
	if err != nil {
		fatalQueryError("account selection using email")
	}
	defer rows.Close()
	if rows.Next() {
		//check password
		var accountID int
		var passhashstr string
		var saltstr string
		rows.Scan(&accountID,&passhashstr,&saltstr)
		salt, _ := hex.DecodeString(saltstr)
		_passhash, _ := hex.DecodeString(passhashstr)
		var passhash [32]byte
		copy(passhash[:],_passhash)
		if passhash == sha256.Sum256( append( []byte(passwd) , salt... ) ) {
			//password match - login the user
			//create session for user
			sessionid := make([]byte, 32)
			_, err := rand.Read(sessionid)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
			sessionidstr := hex.EncodeToString(sessionid)
			//delete any earlier sessions and add this new session
			mysqldb.Exec("delete from sessions where accountid=?",accountID)
			mysqldb.Exec("insert into sessions(sessionid,accountid) values(?,?)",sessionidstr,accountID)
			//send cookie to client
			http.SetCookie(writer , &http.Cookie{Name:"session",Value:sessionidstr,Path:"/"})
			//finally redirect to home page
			http.Redirect(writer, request, "/home/", http.StatusFound)
		} else {
			//wrong password
			http.Redirect(writer, request, "/login/passwderr/", http.StatusFound)
		}
	} else {
		//invalid email - go back to login page with error
		http.Redirect(writer, request, "/login/emailerr/", http.StatusFound)
	}
}

//handler for logouts
func logoutHandler(writer http.ResponseWriter, request *http.Request) {
	//get the "session" cookie provided in the HTTP request
	cookie, err := request.Cookie("session")
	if err == nil {
		//delete the session, then redirect to login page
		mysqldb.Exec("delete from sessions where sessionid=?",cookie.Value)
	}
	http.Redirect(writer, request, "/login/", http.StatusFound)
}
