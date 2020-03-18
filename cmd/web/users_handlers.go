package main

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/asankov/containerizor/pkg/models"
)

func (app *server) handleUserCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.PostFormValue("username")
		if username == "" {
			http.Error(w, "username cannot be empty", http.StatusBadRequest)
			return
		}

		password := r.PostFormValue("password")
		if password == "" {
			http.Error(w, "password cannot be empty", http.StatusBadRequest)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
		if err != nil {
			app.log.Println(err.Error())
			http.Error(w, "internal error", 500)
			return
		}

		userWithSameUsername, err := app.db.GetUserByUsername(username)
		if err != nil || userWithSameUsername != nil {
			app.log.Println(err.Error())
			http.Error(w, "user with same username exists", 400)
			return
		}

		err = app.db.CreateUser(&models.User{
			Username:       username,
			HashedPassword: string(hashedPassword),
		})
		if err != nil {
			app.log.Println(err.Error())
			http.Error(w, "internal error", 500)
			return
		}

		redirectToView(w, "/")
	}
}

func (app *server) handleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.PostFormValue("username")
		if username == "" {
			http.Error(w, "username cannot be empty", http.StatusBadRequest)
			return
		}

		password := r.PostFormValue("password")
		if password == "" {
			http.Error(w, "password cannot be empty", http.StatusBadRequest)
			return
		}

		user, err := app.db.GetUserByUsername(username)
		if err != nil {
			app.log.Println(err.Error())
			http.Error(w, "internal error", 500)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
		if user == nil || err != nil {
			http.Error(w, "username or password are wrong", 400)
			return
		}

		w.Header().Add("Set-Cookie", "token=TODO")
		app.serveTemplate(w, nil, "./ui/html/list.page.tmpl", "./ui/html/base.layout.tmpl")
	}
}

func (app *server) handleLoginView() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		app.serveTemplate(w, nil, "./ui/html/login.page.tmpl", "./ui/html/base.layout.tmpl")
	}
}

func (app *server) handleUserCreateView() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		app.serveTemplate(w, nil, "./ui/html/signup.page.tmpl", "./ui/html/base.layout.tmpl")
	}
}
