package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type loginErr struct {
	UsernameErr string
	PasswordErr string
	Username    string
}

var (
	username = "junaid"
	password = "password"
	store    = sessions.NewCookieStore([]byte("superKey"))
)

func showLoginPage(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("/home/junaid/LoginPage/Templates/login.html")
	if err != nil {
		http.Error(w, "Unable to load template", http.StatusInternalServerError)
		return
	}

	session, err := store.Get(r, "log-session")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if loggedIn, ok := session.Values["logged"].(bool); ok && loggedIn {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	var errInstance loginErr

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Unable to parse the form", http.StatusBadRequest)
			return
		}

		inputUsername := r.FormValue("Username")
		inputPassword := r.FormValue("Password")

		if inputUsername != username {
			errInstance.UsernameErr = "Invalid Username"
		}
		if inputPassword != password {
			errInstance.PasswordErr = "Invalid Password"
		}

		if errInstance.UsernameErr == "" && errInstance.PasswordErr == "" {
			session.Values["logged"] = true
			if err := session.Save(r, w); err != nil {
				http.Error(w, "Unable to save session", http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/home", http.StatusSeeOther)
			return
		}

		errInstance.Username = inputUsername
	}

	if err := temp.Execute(w, errInstance); err != nil {
		http.Error(w, "Unable to execute template", http.StatusInternalServerError)
	}
}

func showHomePage(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "log-session")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if loggedIn, ok := session.Values["logged"].(bool); !ok || !loggedIn {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	temp, err := template.ParseFiles("/home/junaid/LoginPage/Templates/home.html")
	if err != nil {
		http.Error(w, "Unable to load template", http.StatusInternalServerError)
		return
	}

	if err := temp.Execute(w, username); err != nil {
		http.Error(w, "Unable to execute template", http.StatusInternalServerError)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "log-session")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	session.Options.MaxAge = -1
	session.Values["logged"] = false
	if err := session.Save(r, w); err != nil {
		http.Error(w, "Unable to save session", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/login", showLoginPage).Methods("GET", "POST")
	r.HandleFunc("/home", showHomePage).Methods("GET")
	r.HandleFunc("/logout", logoutHandler).Methods("GET")
	http.ListenAndServe(":8080", r)
}
