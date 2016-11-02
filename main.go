package main

import (
	"fmt"
	"os"
	"net/http"
	"html/template"
)

//the template files for the webpages - must be updated as we add pages
var templates = template.Must(template.ParseFiles("templates/login.html",
	"templates/register.html", "templates/home.html", "templates/deanhome.html", 
	"templates/profile.html", "templates/home.css", "templates/course.html",
	"templates/chgroles.html","templates/createCourse.html", "templates/privileges.html",
	"templates/students.html", "templates/roleRedirect.html", "templates/assnCreate.html",
	"templates/assn.html","templates/gradeOverview.html", "templates/grade.html"))

//terminate program if an error occurs while executing a critical query
func fatalQueryError(query string) {
	fmt.Printf("Fatal error while executing query:\n%s\nExiting...\n",query)
	os.Exit(-1)
}

//get account id from cookie
func getAccFromCookie(writer http.ResponseWriter, request *http.Request, redirectToLogin bool) int {
	var accountID int
	//get the "session" cookie provided in the HTTP request
	cookie, err := request.Cookie("session")
	if err != nil {
		//no cookie - the user is visiting for the first time. redirect to login page if requested.
		if redirectToLogin {
			http.Redirect(writer, request, "/login/", http.StatusFound)
		}
		return -1
	}
	//get the account id associated with the session id
	rows , err := mysqldb.Query("select accountid from sessions where sessionid=?",cookie.Value)
	if err != nil {
		fatalQueryError("session id selection")
	}
	defer rows.Close()
	if !rows.Next() {
		//invalid session id - the cookie might be faked. redirect to login page if requested.
		if redirectToLogin {
			http.Redirect(writer, request, "/login/", http.StatusFound)
		}
		return -1
	} else {
		rows.Scan(&accountID)
	}
	return accountID
}

func main() {
	//connect to the MySQL DB
	connect_to_db()
	
	//after main exits, end connection with db
	defer mysqldb.Close()
	
	//delete all sessions
	mysqldb.Exec("delete from sessions")

	//list of handlers - must be updated as we add pages
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/login/",loginHandler)
	http.HandleFunc("/login-submit/",loginSubmitHandler)
	http.HandleFunc("/logout/",logoutHandler)
	http.HandleFunc("/register/",registerHandler)
	http.HandleFunc("/register-submit/",registerSubmitHandler)
	http.HandleFunc("/home/",homeHandler)
	http.HandleFunc("/deanhome/",deanhomeHandler)
	http.HandleFunc("/priv-esc/", privilegeHandler)
	http.HandleFunc("/profile/",profileHandler)
	http.HandleFunc("/course/",courseHandler)
	http.HandleFunc("/chgroles/",courseRolesHandler)
	http.HandleFunc("/roleAssn/",roleAssnHandler)
	http.HandleFunc("/createCourse/",createCourseHandler)
	http.HandleFunc("/course-submit/",createCourseSubmitHandler)
	http.HandleFunc("/privileges/",privilegeHandler)
	http.HandleFunc("/privilegeEscalator/", privilegeEscHandler)
	http.HandleFunc("/students/", studentHandler)
	http.HandleFunc("/assnCreate/", assnCreateHandler)
	http.HandleFunc("/registerAssn/", assnRegisterHandler)
	http.HandleFunc("/download/",downloadHandler)
	http.HandleFunc("/assn/",assignmentHandler)
	http.HandleFunc("/comment-submit/",commentSubmitHandler)
	http.HandleFunc("/answer-submit/",commentSubmitHandler)
	//http.HandleFunc("/createAssn/",createAssignmentHandler)
	http.HandleFunc("/grade/",gradeHandler)
	http.HandleFunc("/grade-submit/",gradeSubmitHandler)

	//Don't know how to redirect css files. Not required at the moment
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("http/css"))))
	
	//start server
	http.ListenAndServe(":8080", nil)
}
