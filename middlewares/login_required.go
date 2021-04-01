package middlewares

import (
	"gallerio/models"
	"gallerio/utils/context"
	"net/http"
)

type AssignUser struct {
	models.UserService
}

func (mw *AssignUser) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFunc(next.ServeHTTP)
}

func (mw *AssignUser) ApplyFunc(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie("remember_token")
		if err != nil {
			next(w, req)
			return
		}
		
		user, err := mw.UserService.ByRememberToken(cookie.Value)
		if err != nil {
			next(w, req)
			return
		}
		ctx := req.Context()
		ctx = context.WithUser(ctx, user)
		req = req.WithContext(ctx)
		
		next(w, req)
	}
}

type LoginRequired struct {
	models.UserService
}

func (mw *LoginRequired) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFunc(next.ServeHTTP)
}

func (mw *LoginRequired) ApplyFunc(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		user := context.User(req.Context())
		if user == nil {
			http.Redirect(w, req, "/signin", http.StatusSeeOther)
			return
		}
		next(w, req)
	}
}
