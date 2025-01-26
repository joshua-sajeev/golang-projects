package main

import (
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/joho/godotenv"
)

var (
	database     *sql.DB
	hashKey      = securecookie.GenerateRandomKey(32)
	secureCookie = securecookie.New(hashKey, nil)
)

type JSONResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Data    map[string]any `json:"data,omitempty"`
	Errors  []string       `json:"errors,omitempty"`
}

type Page struct {
	ID         int           `json:"id"`
	Title      string        `json:"title"`
	RawContent string        `json:"raw_content"`
	Content    template.HTML `json:"content"`
	Date       string        `json:"date"`
	Guid       string        `json:"guid"`
	Comments   []Comment     `json:"comments,omitempty"`
}

type Comment struct {
	ID     int       `json:"id"`
	PageID int       `json:"page_id"`
	Guid   string    `json:"comment_guid"`
	Name   string    `json:"comment_name"`
	Email  string    `json:"email"`
	Text   string    `json:"text"`
	Date   time.Time `json:"comment_date"`
}

func loadEnvVars() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}
func (p Page) TruncatedText() string {
	chars := 0
	for i := range p.Content {
		chars++
		if chars > 150 {
			return string(p.Content[:i]) + "..."
		}
	}
	return string(p.Content)
}

func sanitizeHTML(raw string) template.HTML {
	return template.HTML(raw)
}

func (p *Page) LoadComments() error {
	rows, err := database.Query(
		`SELECT id, comment_guid, comment_name, comment_email, 
         comment_text, comment_date 
         FROM comments 
         WHERE page_id = ? 
         ORDER BY comment_date DESC`,
		p.ID,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var c Comment
		err := rows.Scan(
			&c.ID,
			&c.Guid,
			&c.Name,
			&c.Email,
			&c.Text,
			&c.Date,
		)
		if err != nil {
			return err
		}
		c.PageID = p.ID
		p.Comments = append(p.Comments, c)
	}
	return rows.Err()
}

func APIPages(w http.ResponseWriter, r *http.Request) {
	var pages []Page
	rows, err := database.Query(
		"SELECT id, page_title, page_content, page_date, page_guid FROM pages ORDER BY page_date DESC",
	)
	if err != nil {
		sendJSONError(w, "Database error", http.StatusInternalServerError)
		log.Printf("Database error in APIPages: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var p Page
		err := rows.Scan(&p.ID, &p.Title, &p.RawContent, &p.Date, &p.Guid)
		if err != nil {
			continue
		}
		p.Content = sanitizeHTML(p.RawContent)
		if err := p.LoadComments(); err != nil {
			log.Printf("Error loading comments for page %s: %v", p.Guid, err)
		}
		pages = append(pages, p)
	}

	sendJSONResponse(w, JSONResponse{
		Success: true,
		Data:    map[string]any{"pages": pages},
	})
}

func APIPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pageGUID := vars["guid"]

	var p Page
	err := database.QueryRow(
		"SELECT id, page_title, page_content, page_date, page_guid FROM pages WHERE page_guid = ?",
		pageGUID,
	).Scan(&p.ID, &p.Title, &p.RawContent, &p.Date, &p.Guid)

	if err != nil {
		if err == sql.ErrNoRows {
			sendJSONError(w, "Page not found", http.StatusNotFound)
		} else {
			sendJSONError(w, "Database error", http.StatusInternalServerError)
			log.Printf("Database error in APIPage: %v", err)
		}
		return
	}

	p.Content = sanitizeHTML(p.RawContent)
	if err := p.LoadComments(); err != nil {
		log.Printf("Error loading comments for page %s: %v", p.Guid, err)
	}

	sendJSONResponse(w, JSONResponse{
		Success: true,
		Data:    map[string]any{"page": p},
	})
}

func APICommentPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		sendJSONError(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Get page_id from page_guid
	pageGuid := r.FormValue("guid")
	var pageID int
	err := database.QueryRow(
		"SELECT id FROM pages WHERE page_guid = ?",
		pageGuid,
	).Scan(&pageID)
	if err != nil {
		sendJSONError(w, "Invalid page GUID", http.StatusBadRequest)
		log.Printf("Error getting page_id: %v", err)
		return
	}

	comment := Comment{
		PageID: pageID,
		Guid:   generateGUID(),
		Name:   r.FormValue("name"),
		Email:  r.FormValue("email"),
		Text:   r.FormValue("comments"),
	}

	// Validate input
	var errors []string
	if len(comment.Name) < 2 {
		errors = append(errors, "Name must be at least 2 characters long")
	}
	if len(comment.Email) < 5 || !strings.Contains(comment.Email, "@") {
		errors = append(errors, "Invalid email address")
	}
	if len(comment.Text) < 1 {
		errors = append(errors, "Comment text is required")
	}

	if len(errors) > 0 {
		sendJSONResponse(w, JSONResponse{
			Success: false,
			Message: "Validation failed",
			Errors:  errors,
		})
		return
	}

	// Insert comment
	res, err := database.Exec(
		`INSERT INTO comments 
         (page_id, comment_guid, comment_name, comment_email, comment_text, comment_date) 
         VALUES (?, ?, ?, ?, ?, NOW())`,
		comment.PageID,
		comment.Guid,
		comment.Name,
		comment.Email,
		comment.Text,
	)

	if err != nil {
		sendJSONError(w, "Error saving comment", http.StatusInternalServerError)
		log.Printf("Database error in APICommentPost: %v", err)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		sendJSONError(w, "Error retrieving comment ID", http.StatusInternalServerError)
		return
	}

	comment.ID = int(id)
	comment.Date = time.Now()

	sendJSONResponse(w, JSONResponse{
		Success: true,
		Message: "Comment added successfully",
		Data:    map[string]any{"comment": comment},
	})
}

func ServeIndex(w http.ResponseWriter, r *http.Request) {
	var pages []Page
	rows, err := database.Query(
		"SELECT id, page_title, page_content, page_date, page_guid FROM pages ORDER BY page_date DESC",
	)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		log.Printf("Database error in ServeIndex: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var p Page
		err := rows.Scan(&p.ID, &p.Title, &p.RawContent, &p.Date, &p.Guid)
		if err != nil {
			continue
		}
		p.Content = sanitizeHTML(p.RawContent)
		pages = append(pages, p)
	}

	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		log.Printf("Template parsing error: %v", err)
		return
	}

	if err := t.Execute(w, pages); err != nil {
		log.Printf("Template execution error: %v", err)
	}
}

func ServePage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pageGUID := vars["guid"]

	var p Page
	err := database.QueryRow(
		"SELECT id, page_title, page_content, page_date, page_guid FROM pages WHERE page_guid = ?",
		pageGUID,
	).Scan(&p.ID, &p.Title, &p.RawContent, &p.Date, &p.Guid)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Page not found", http.StatusNotFound)
		} else {
			http.Error(w, "Server error", http.StatusInternalServerError)
			log.Printf("Database error in ServePage: %v", err)
		}
		return
	}

	p.Content = sanitizeHTML(p.RawContent)
	if err := p.LoadComments(); err != nil {
		log.Printf("Error loading comments for page %s: %v", p.Guid, err)
	}

	t, err := template.ParseFiles("templates/blog.html")
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		log.Printf("Template parsing error: %v", err)
		return
	}

	if err := t.Execute(w, p); err != nil {
		log.Printf("Template execution error: %v", err)
	}
}

func generateGUID() string {
	return fmt.Sprintf("%d-%x", time.Now().UnixNano(),
		securecookie.GenerateRandomKey(6))
}

func sendJSONResponse(w http.ResponseWriter, resp JSONResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

func sendJSONError(w http.ResponseWriter, message string, code int) {
	w.WriteHeader(code)
	sendJSONResponse(w, JSONResponse{
		Success: false,
		Message: message,
	})
}

func main() {
	loadEnvVars()

	// Read configuration from environment variables
	DBHost := os.Getenv("DB_HOST")
	DBPort := os.Getenv("DB_PORT")
	DBUser := os.Getenv("DB_USER")
	DBPass := os.Getenv("DB_PASS")
	DBDbase := os.Getenv("DB_NAME")
	HTTPPort := os.Getenv("HTTP_PORT")
	HTTPSPort := os.Getenv("HTTPS_PORT")
	KeyFile := os.Getenv("KEY_FILE")
	CertFile := os.Getenv("CERT_FILE")
	// Initialize database connection
	dbConn := fmt.Sprintf("%s:%s@tcp(%s%s)/%s?parseTime=true", DBUser, DBPass, DBHost, DBPort, DBDbase)
	db, err := sql.Open("mysql", dbConn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	database = db
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Configure router
	router := mux.NewRouter()

	// API routes
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/pages", APIPages).Methods("GET")
	api.HandleFunc("/page/{guid}", APIPage).Methods("GET")
	api.HandleFunc("/comments", APICommentPost).Methods("POST")

	// Web routes
	router.HandleFunc("/page/{guid:[0-9a-zA-Z-]+}", ServePage).Methods("GET")
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/home", http.StatusMovedPermanently)
	})
	router.HandleFunc("/home", ServeIndex)

	// Static files
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir("static"))))

	// TLS configuration
	tlsConfig := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256, // Required for HTTP/2
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,   // Required for HTTP/2
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		},
	}

	// Create servers
	httpServer := &http.Server{
		Addr:         HTTPPort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	httpsServer := &http.Server{
		Addr:         HTTPSPort,
		Handler:      router,
		TLSConfig:    tlsConfig,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start servers
	go func() {
		log.Printf("Starting HTTP server on %s", HTTPPort)
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	log.Printf("Starting HTTPS server on %s", HTTPSPort)
	if err := httpsServer.ListenAndServeTLS(CertFile, KeyFile); err != http.ErrServerClosed {
		log.Fatalf("HTTPS server failed: %v", err)
	}
}
