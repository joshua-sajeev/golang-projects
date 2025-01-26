package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

const (
	DBHost  = "127.0.0.1"
	DBPort  = ":3306"
	DBUser  = "root"
	DBPass  = "joshua"
	DBDbase = "webappbook"
	PORT    = ":8080"
)

var database *sql.DB

type Page struct {
	Title      string
	RawContent string
	Content    template.HTML
	Date       string
	GUID       string
}

func (p Page) TruncatedText() string {
	chars := 0
	for i := range p.Content {
		chars++
		if chars > 150 {
			return string(p.Content[:i] + `...`)
		}
	}
	return string(p.Content)
}

func ServeIndex(w http.ResponseWriter, r *http.Request) {
	var Pages = []Page{}
	pages, err := database.Query("SELECT page_title,page_content,page_date,page_guid FROM pages ORDER BY ? DESC", "page_date")
	if err != nil {
		fmt.Println(w, err.Error())
	}

	defer pages.Close()
	for pages.Next() {
		thisPage := Page{}
		pages.Scan(&thisPage.Title, &thisPage.RawContent, &thisPage.Date, &thisPage.GUID)
		thisPage.Content = template.HTML(thisPage.RawContent)
		Pages = append(Pages, thisPage)
	}

	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, Pages)
}

func ServePage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pageGUID := vars["guid"]

	thisPage := Page{}
	fmt.Println(pageGUID)
	err := database.QueryRow("SELECT page_title, page_content, page_date FROM pages WHERE page_guid=?", pageGUID).
		Scan(&thisPage.Title, &thisPage.RawContent, &thisPage.Date)
	thisPage.Content = template.HTML(thisPage.RawContent)

	if err != nil {
		http.Error(w, http.StatusText(404), http.StatusNotFound)
		log.Println("Couldn't get page!")
		return
	}

	// Html using templates
	t, _ := template.ParseFiles("templates/blog.html")
	t.Execute(w, thisPage)
}

func RedrIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/home", 301)
}

func main() {
	dbConn := fmt.Sprintf("%s:%s@/%s", DBUser, DBPass, DBDbase)
	fmt.Println("Connecting to database:", dbConn)

	db, err := sql.Open("mysql", dbConn)
	if err != nil {
		log.Println("Couldn't connect to database:", DBDbase)
		log.Println("Error:", err.Error())
		return
	}

	database = db

	routes := mux.NewRouter()
	// routes.HandleFunc("/page/{id:[0-9]+}", ServePage)
	routes.HandleFunc("/page/{guid:[0-9a-zA\\-]+}", ServePage)
	routes.HandleFunc("/", RedrIndex)
	routes.HandleFunc("/home", ServeIndex)
	http.Handle("/", routes)
	log.Println("Server starting on", PORT)
	if err := http.ListenAndServe(PORT, nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
