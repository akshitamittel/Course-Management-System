use cms;

drop table if exists subfile cascade;
drop table if exists qfile cascade;
drop table if exists comment cascade;
drop table if exists substring cascade;
drop table if exists gracedays cascade;
drop table if exists file cascade;
drop table if exists grade cascade;
drop table if exists assignment cascade;
drop table if exists role cascade;
drop table if exists course cascade;
drop table if exists sessions cascade;
drop table if exists account cascade;

create table account
(accountid int primary key auto_increment,
firstname varchar(64) not null,
lastname varchar(64) not null,
email varchar(64) not null,
passhash char(64) not null,
salt char(32) not null,
rollno varchar(32),
privilegelevel int not null,
activationkey char(64));

create table sessions
(sessionid char(64) not null,
accountid int not null , foreign key (accountid) references account(accountid));

create table course
(courseid int primary key auto_increment,
coursename varchar(64) not null,
coursecode varchar(16) not null,
creatorid int not null , foreign key (creatorid) references account(accountid),
startdate date not null,
enddate date not null,
maxgracedays int not null);

create table role
(accountid int not null , foreign key (accountid) references account(accountid),
courseid int not null , foreign key (courseid) references course(courseid),
role int not null);

create table assignment
(assignmentid int primary key auto_increment,
courseid int not null , foreign key (courseid) references course(courseid),
creatorid int not null , foreign key (creatorid) references account(accountid),
creationtime datetime not null,
duetime datetime,
maxsubmittime datetime,
titlestring text not null,
submissiontype varchar(4) not null,
maxmarks int not null);

create table grade
(assignmentid int not null , foreign key (assignmentid) references assignment(assignmentid),
submittorid int not null , foreign key (submittorid) references account(accountid),
givenmarks int,
gracedaysused int not null);

create table file
(fileid int primary key auto_increment,
realfname varchar(64) not null,
storedfname varchar(96) not null);

create table subfile
(fileid int not null , foreign key (fileid) references file(fileid),
submittorid int not null , foreign key (submittorid) references account(accountid),
assignmentid int not null , foreign key (assignmentid) references assignment(assignmentid),
submittime datetime not null default CURRENT_TIMESTAMP);

create table qfile
(fileid int not null , foreign key (fileid) references file(fileid),
assignmentid int not null , foreign key (assignmentid) references assignment(assignmentid));

create table comment
(commentstring text not null,
commentorid int not null , foreign key (commentorid) references account(accountid),
assignmentid int not null , foreign key (assignmentid) references assignment(assignmentid),
commenttime datetime not null);

create table substring
(answerstring text not null,
submittorid int not null , foreign key (submittorid) references account(accountid),
assignmentid int not null , foreign key (assignmentid) references assignment(assignmentid),
submittime datetime not null);

create table gracedays
(courseid int not null , foreign key (courseid) references course(courseid),
studentid int not null , foreign key (studentid) references account(accountid),
availabledays int not null);
