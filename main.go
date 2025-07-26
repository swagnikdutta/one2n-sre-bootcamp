package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

const (
	studentIdKey contextKey = "studentId"

	createTableSyntax = `create table if not exists students (
		id integer primary key autoincrement,
		name text not null,
		age integer 
	)`
	insertSyntax        = `insert into students (id, name, age) values (?, ?, ?)`
	listSyntax          = `select * from students`
	selectStudentSyntax = `select * from students where id = ?`
	deleteSyntax        = `delete from students where id = ?`
	updateSyntax        = `update students set name = ?, age = ? where id = ?`
)

func NewRequestMultiplexer(server *Server) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/students", server.listStudents)
	mux.HandleFunc("/students/add", server.addStudent)
	mux.HandleFunc("/students/{id}", server.studentHandler)
	return mux
}

func main() {
	db, err := sql.Open("sqlite3", "students.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	server := NewServer(db)
	httpServer := &http.Server{
		Addr:    ":8000",
		Handler: NewRequestMultiplexer(server),
	}

	server.initDB()

	err = httpServer.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
