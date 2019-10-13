package http

import (
	"net/http"
	"strconv"
	app "useritem"
	"useritem/context"
)

// Auth is authentication middleware
type AuthMw struct {
	userRepo app.UserRepo
}

// UserViaSession retrieves a user from session
// and put it into request context
func (a *AuthMw) UserViaSession(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionStr, err := r.Cookie("session")
		if err != nil {
			// No user session found, move on
			next.ServeHTTP(w, r)
			return
		}

		// Convert to integer
		session, err := strconv.Atoi(sessionStr.Value)
		if err != nil {
			// Session is not integer, move on
			next.ServeHTTP(w, r)
			return
		}

		user, err := a.userRepo.ByToken(session)
		if err != nil {
			// No user found, move on
			next.ServeHTTP(w, r)
			return
		}
		r = r.WithContext(context.WithUser(r.Context(), user))
		next.ServeHTTP(w, r)
	}
}

// RequireUser requires a user from context
// if no user found, redirect to sign in
func (a *AuthMw) RequireUser(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmp := context.User(r.Context())
		if tmp == nil {
			// No user found
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	}
}
