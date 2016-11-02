package main

import (
	"net/http"
)

//handler for the 'index' page i.e. empty URL
func indexHandler(writer http.ResponseWriter, request *http.Request) {
	if getAccFromCookie(writer, request, true) >= 0 {
		//a valid account id is associated with this session id, so go to homepage
		http.Redirect(writer, request, "/home/", http.StatusFound)
	}
}
