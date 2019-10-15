package http

import (
	"net/http"
)

// AuthMw is authentication middleware
type AuthMw interface {
	SetUser(next http.Handler) http.HandlerFunc
	RequireUser(next http.Handler) http.HandlerFunc
}
