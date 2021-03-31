package middlewares

import (
	"gallerio/models"
	"gallerio/utils/context"
	"net/http"
)

type LoginRequired struct {
	models.UserService
}

func (mw *LoginRequired) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFunc(next.ServeHTTP)
}

// Interrupts requests and gets the remember_token from cookies
// Then checks if user exists in database
// If exists, sets the user in the context
// Otherwise redirects to login page
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

		// Adds current user to context
		ctx := req.Context()
		ctx = context.WithUser(ctx, user)
		req = req.WithContext(ctx)

		// Continues Request
		next(w, req)
	}
}