package main

import (
	"fmt"
	"net/http"
	"time"
)

type Login struct {
	HashedPassword string
	SessionToken   string
	CSRFToken      string
}

var users = map[string]Login{}

func main() {
	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.HandleFunc("/protected", protected)
	http.HandleFunc("/logout", logout)

	http.ListenAndServe(":8080", nil)
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if len(username) < 8 || len(password) < 8 {
		http.Error(w, "Invalid username/password", http.StatusNotAcceptable)
		return
	}

	if _, ok := users[username]; ok {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	users[username] = Login{
		HashedPassword: hashedPassword,
	}
	fmt.Fprintf(w, "User registered successfully!")
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	user, ok := users[username]
	if !ok || !checkPasswordHash(password, user.HashedPassword) {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Generate new session token and CSRF token
	sessionToken := generateToken(32)
	csrfToken := generateToken(32)

	// Set secure session cookie
	sessionCookie := &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true, // Only send over HTTPS
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	}
	http.SetCookie(w, sessionCookie)

	// Set CSRF cookie
	csrfCookie := &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: false, // Accessible by JavaScript
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	}
	http.SetCookie(w, csrfCookie)

	// Update user's session information
	user.SessionToken = sessionToken
	user.CSRFToken = csrfToken
	users[username] = user

	fmt.Fprintln(w, "Login successful")
}

func protected(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := Authorize(r); err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	username := r.FormValue("username")
	fmt.Fprintf(w, "CSRF validation successful! Welcome %s", username)
}

func logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Clear session cookie
	sessionCookie := &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		MaxAge:   -1,
	}
	http.SetCookie(w, sessionCookie)

	// Clear CSRF cookie
	csrfCookie := &http.Cookie{
		Name:     "csrf_token",
		Value:    "",
		Expires:  time.Now().Add(-24 * time.Hour),
		HttpOnly: false,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		MaxAge:   -1,
	}
	http.SetCookie(w, csrfCookie)

	username := r.FormValue("username")
	if user, ok := users[username]; ok {
		user.SessionToken = ""
		user.CSRFToken = ""
		users[username] = user
	}

	fmt.Fprintf(w, "Logged out successfully")
}
