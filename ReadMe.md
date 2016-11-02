#DBMS Project: College Management System in Golang.

##Authors:
* Akshita Mittel (cs13b1040@iith.ac.in)
* Keyur Joshi (cs13b1011@iith.ac.in)
* CH Sreekar Reddy (cs13b1008@iith.ac.in)
* Arjun V Anand (cs13b1041@iith.ac.in)

##E-R digram of the Database system
![](https://github.com/akshitamittel/Course-Management-System/blob/master/diagram.png)

##Instructions to run the Database
1. You will need mysql and go-sql-driver. 
  * For getting the second, you will need to configure your GOPATH and GOBIN directories.
  * gopath should contain 3 folders : bin, pkg, src. 
  * gobin should be the bin folder in gopath.
2. In mysql, create database called 'cms'. 
  * then create user called 'cms' with password 'cms' and grant it all permissions on database 'cms'. 
  * To run mysql do: 
  ```
  mysql -u <username> -p
  ```
  then press enter, then enter password. 
  * to execute from file, do:
  ```
  mysql -u <username> -p < <filename>
  ```
  then press enter, then enter password.
3. execute make-tables.sql using above instructions
4. compile using 
  ```
  go install cms
  ```
  (should work if gopath and gobin are configured)
5. run copy html to bin script
6. cd to the bin directory and run as: ```./cms```
7. you can now run the site. access it at localhost:8080 address.
