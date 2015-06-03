package chatserver

import (
	"database/sql"
	//database driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/matter123/chat/config"
	"github.com/matter123/chat/token"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

func setStatus(w http.ResponseWriter, status string, url string) {
	secure := false
	if config.Settings().SSLSettings != nil {
		secure = true
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "status",
		Value:    status,
		Path:     "/",
		Secure:   secure,
		Expires:  time.Now().Add(1 * time.Second),
		HttpOnly: false,
	})
	w.Header().Set("Location", url)
	w.WriteHeader(303)
}

//TODO move DB code out of login.go
var db = setupDB()

func setupDB() *sql.DB {
	dsn := config.Settings().MYSQLSettings.User + ":" + config.Settings().MYSQLSettings.Pass + "@tcp(" +
		config.Settings().MYSQLSettings.Host + ")/" + config.Settings().MYSQLSettings.DB
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return db
}

//LoginHandle is the handle to /login, LoginHandle checks for token cookies, and check validity
//If no tokens are provided, LoginHandle takes a form encoded user and password, and logs in,
//setting the token cookie to the users session token
func LoginHandle(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err == nil && cookie != nil {
		//read token cookie, check if still valid
		//if token is still valid redirect to /chat.html
		//else drop the cookie, and redirect to index.html
		if token.Valid(cookie.Value) {
			w.Header().Set("Location", "chat.html")
			w.WriteHeader(303)
			return
		}
		//token is invalid, delete it and redir to index.html
		secure := false
		if config.Settings().SSLSettings != nil {
			secure = true
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    "",
			Path:     "/",
			Expires:  time.Now().Add(-1 * time.Hour),
			Secure:   secure,
			HttpOnly: true,
		})
		w.Header().Set("Location", "chat.html")
		w.WriteHeader(303)
		return
	}
	user := r.FormValue("user")
	pass := r.FormValue("pass")
	if user == "" || pass == "" {
		//username and/or password not provided
		setStatus(w, "username and/or password not provided", "index.html")
		return
	}
	//get hash from DB and compare to provided pass
	row := db.QueryRow("SELECT hash FROM users WHERE username = ?", user)
	hash := ""
	err = row.Scan(&hash)

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	if err != nil {
		setStatus(w, "username and/or password incorrect", "index.html")
		return
	}
	//login is valid
	//set cookie and redirect
	secure := false
	if config.Settings().SSLSettings != nil {
		secure = true
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token.Token(user),
		Path:     "/",
		Secure:   secure,
		HttpOnly: true,
	})
	w.Header().Set("Location", "chat.html")
	w.WriteHeader(303)
}

//SignupHandle handles /signup, takes a form encoded username and password, checks the
//database for an existing user, and sets the session token cookie, and redirects to the chat
func SignupHandle(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("user")
	pass := r.FormValue("pass")
	if user == "" || pass == "" {
		//username and/or password not provided
		setStatus(w, "username and/or password not provided", "signup.html")
		return
	}
	trash := ""
	err := db.QueryRow("SELECT username FROM users WHERE username = ?", user).Scan(&trash)
	if err != sql.ErrNoRows {
		setStatus(w, "username already exists", "signup.html")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(pass), config.Settings().PasswordSettings.Cost)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("INSERT INTO users (username, hash) VALUES (?, ?)", user, string(hash))
	if err != nil {
		log.Fatal(err)
	}
	//set cookie and redirect
	secure := false
	if config.Settings().SSLSettings != nil {
		secure = true
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token.Token(user),
		Path:     "/",
		Secure:   secure,
		HttpOnly: true,
	})
	w.Header().Set("Location", "chat.html")
	w.WriteHeader(303)
}
