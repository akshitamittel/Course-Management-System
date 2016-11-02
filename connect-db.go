package main

import (
	"fmt"
	"os"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

//this variable is visible across the package
var mysqldb *sql.DB

//try to connect to the mysql database
func connect_to_db() {
	//the format for the second argument is "<username>:<password>@/<dbname>"
	//we will have to use the same username, password name, and database name, or else we will have trouble with git
	var err error
	mysqldb, err = sql.Open("mysql", "cms:cms@/cms")
	if err != nil {
		//error in opening database - the program cannot proceed any further, so just exit and return -1
		fmt.Printf("Fatal error while connecting to database. Exiting...\n")
		os.Exit(-1)
	}
}
