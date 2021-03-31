package middlewares

import (
	"gallerio/accounts"
	"log"
	"net/http"
)

type LoginRequired struct {
	accounts.UserService
}

func (mw *LoginRequired) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFunc(next.ServeHTTP)
}

func (mw *LoginRequired) ApplyFunc(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie("remember_token")
		if err != nil {
			http.Redirect(w, req, "/signin", http.StatusSeeOther)
			return
		}

		user, err := mw.UserService.ByRememberToken(cookie.Value)
		if err != nil {
			http.Redirect(w, req, "/signin", http.StatusSeeOther)
			return
		}
		log.Println(user)

		next(w, req)
	}
}