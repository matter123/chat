package main

import (
	"net/http"
	//"database/sql"

)

func set_status(w http.ResponseWriter, status string, url string) {
	secure := false
	if config.SSLSettings != nil {
		secure = true
	}
	http.SetCookie(w, &http.Cookie{
		Name: "status",
		Value: status,
		Path: "/",
		Secure: secure,
		HttpOnly: true,
	})
	w.Header().Set("Location", url)
	w.WriteHeader(303)
}

func login(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err == nil && cookie != nil {
		//read token cookie, check if still valid
		//if token is still valid redirect to /chat.html
		//else drop the cookie, and redirect to index.html

		//for now assume cookie is always valid
		w.Header().Set("Location", "chat.html")
		return;
	}
	user := r.FormValue("user")
	pass := r.FormValue("pass")
	if user == "" || pass == "" {
		//username and/or password not provided
		set_status(w, "username and/or password not provided", "index.html")
		return
	}
}

func signup(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("user")
	pass := r.FormValue("pass")
	if user == "" || pass == "" {
		//username and/or password not provided
		set_status(w, "username and/or password not provided", "index.html")
		return
	}
}
