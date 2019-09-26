package middleware

import (
	"context"
	"net/http"
	"strconv"
	app "useritem"
)

type Auth struct {
	UserRepo app.UserRepo
}

// UserViaSession retrieves a user from session
// and put it into request context
func (a *Auth) UserViaSession(next http.Handler) http.HandlerFunc {
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

		user, err := a.UserRepo.ByToken(session)
		if err != nil {
			// No user found, move on
			next.ServeHTTP(w, r)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), "user", user))
		next.ServeHTTP(w, r)
	}
}

// RequireUser requires a user from context
// if no user found, redirect to sign in
func (a *Auth) RequireUser(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmp := r.Context().Value("user")
		if tmp == nil {
			// No user found
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		}
		if _, ok := tmp.(*app.User); !ok {
			// Value from context is not an user
			http.Redirect(w, r, "signin", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	}
}
