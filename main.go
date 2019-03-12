package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

type address_book struct {
	Id    int
	Fname string
	Lname string
	Email string
	Phone string
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "tester"
	dbPass := "password"
	dbName := "pizza_hut"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

var tmpl = template.Must(template.ParseGlob("form/*"))

func Index(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT * From address_book ORDER BY id DESC")
	if err != nil {
		panic(err.Error())
	}
	add := address_book{}
	res := []address_book{}
	for selDB.Next() {
		var id int
		var fname, lname, email, phone string
		err = selDB.Scan(&id, &fname, &lname, &email, &phone)
		if err != nil {
			panic(err.Error())
		}
		add.Id = id
		add.Fname = fname
		add.Lname = lname
		add.Email = email
		add.Phone = phone
		res = append(res, add)
	}
	tmpl.ExecuteTemplate(w, "Index", res)
	defer db.Close()
}

func Show(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM address_book WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	add := address_book{}
	for selDB.Next() {
		var id int
		var fname, lname, email, phone string
		err = selDB.Scan(&id, &fname, &lname, &email, &phone)
		if err != nil {
			panic(err.Error())
		}
		add.Id = id
		add.Fname = fname
		add.Lname = lname
		add.Email = email
		add.Phone = phone
	}
	tmpl.ExecuteTemplate(w, "Show", add)
	defer db.Close()
}

func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM address_book WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	add := address_book{}
	for selDB.Next() {
		var id int
		var fname, lname, email, phone string
		err = selDB.Scan(&id, &fname, &lname, &email, &phone)
		if err != nil {
			panic(err.Error())
		}
		add.Id = id
		add.Fname = fname
		add.Lname = lname
		add.Email = email
		add.Phone = phone
	}
	tmpl.ExecuteTemplate(w, "Edit", add)
	defer db.Close()
}

func Insert(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		fname := r.FormValue("fname")
		lname := r.FormValue("lname")
		email := r.FormValue("email")
		phone := r.FormValue("phone")
		insForm, err := db.Prepare("INSERT INTO address_book(fname, lname, email, phone) VALUES(?,?,?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(fname, lname, email, phone)
		log.Println("INSERT: Fname: " + fname + " | Lname: " + lname + " | Email: " + email + " | Phone: " + phone)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Update(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		fname := r.FormValue("fname")
		lname := r.FormValue("lname")
		email := r.FormValue("email")
		phone := r.FormValue("phone")
		id := r.FormValue("uid")
		insForm, err := db.Prepare("UPDATE address_book SET fname=?, lname=?, email=?, phone=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(fname, lname, email, phone, id)
		log.Println("UPDATE: Fname: " + fname + " | Lname: " + lname + " | Email: " + email + " | Phone: " + phone)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	add := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM address_book WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(add)
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func main() {
	log.Println("Server started on:  http://localhost:8080")
	http.HandleFunc("/", Index)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)
	http.ListenAndServe(":8080", nil)
}
