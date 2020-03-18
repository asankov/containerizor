package main

import "net/http"

func (srv *server) requireLogin(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")
		if err != nil {
			redirectToView(w, "/users/login")
			return
		}

		username, err := srv.auth.UserFromToken(c.Value)
		if err != nil {
			redirectToView(w, "/users/login")
			return
		}

		r.Header.Set("user", username)
		srv.log.Println("user authenticated ", username)

		f(w, r)
	}
}
