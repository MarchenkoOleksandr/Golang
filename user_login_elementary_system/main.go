package main

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"io"
	"log"
	"net/http"
)

type Message struct {
	Color string
	Text  string
}

var database *sql.DB

func main() {
	db, err := sql.Open("mysql", "root:password@/userLogin")

	if err != nil {
		log.Println(err)
	}

	database = db
	defer db.Close()

	http.HandleFunc("/", handler)

	fmt.Println("\nServer is listening on port 8181...")
	_ = http.ListenAndServe(":8181", nil)
}

func handler(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		log.Println(err)
	}

	if request.FormValue("submit") == "log" {
		login(request, writer)
	} else if request.FormValue("submit") == "reg" {
		registration(request, writer)
	} else {
		getRequest(writer)
	}
}

func login(request *http.Request, writer http.ResponseWriter) {
	rows, err := database.Query(
		"SELECT login, passwordHash FROM Users WHERE login = ? AND passwordHash = ?",
		request.FormValue("login"), hash(request.FormValue("pass")))

	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	if rows.Next() == true {
		openWithSuccess(&writer, "Login is successful!")
	} else {
		openWithError(&writer, "Login failed!")
	}
}

func registration(request *http.Request, writer http.ResponseWriter) {
	_, err := database.Exec("INSERT INTO Users (login, passwordHash) values (?, ?)",
		request.FormValue("login"), hash(request.FormValue("pass")))

	if err != nil {
		openWithError(&writer, "Error: a user with this login already exists!")
	} else {
		openWithSuccess(&writer, "Registration is successful!")
	}
}

func getRequest(writer http.ResponseWriter) {
	openWithSuccess(&writer, "")
}

func openWithSuccess(writer *http.ResponseWriter, msg string) {
	message := Message{"success", msg}
	tmpl, _ := template.ParseFiles("templates/index.html")
	_ = tmpl.Execute(*writer, message)
}

func openWithError(writer *http.ResponseWriter, msg string) {
	message := Message{"warning", msg}
	tmpl, _ := template.ParseFiles("templates/index.html")
	_ = tmpl.Execute(*writer, message)
}

func hash(str string) string {
	hashCode := md5.New()
	_, _ = io.WriteString(hashCode, str)
	return fmt.Sprintf("%x", hashCode.Sum(nil))
}
