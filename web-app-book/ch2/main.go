package main

import (
	"database/sql"
	"fmt"
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
	Title   string
	Content string
	Date    string
}

func ServePage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pageGUID := vars["guid"]

	thisPage := Page{}
	fmt.Println(pageGUID)
	err := database.QueryRow("SELECT page_title, page_content, page_date FROM pages WHERE page_guid=?", pageGUID).
		Scan(&thisPage.Title, &thisPage.Content, &thisPage.Date)

	if err != nil {
		http.Error(w, http.StatusText(404), http.StatusNotFound)
		log.Println("Couldn't get page!")
		return
	}
	// Prepare the HTML response
	html := fmt.Sprintf(`
        <html>
            <head><title>%s</title></head>
            <body>
                <h1>%s</h1>
                <div>%s</div>
            </body>
        </html>`,
		thisPage.Title, thisPage.Title, thisPage.Content,
	)

	// Write the HTML response
	fmt.Fprintln(w, html)
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

	http.Handle("/", routes)
	log.Println("Server starting on", PORT)
	if err := http.ListenAndServe(PORT, nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
