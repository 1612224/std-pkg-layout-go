package http

import (
	"net/http"
	app "useritem"

	"github.com/gorilla/mux"
)

// NewServer returns a server that handles both HTML and JSON
func NewServer(userRepo app.UserRepo, itemRepo app.ItemRepo) http.Handler {
	html := HTMLServer(userRepo, itemRepo)
	json := JSONServer(userRepo, itemRepo)
	mux := http.NewServeMux()
	mux.Handle("/", html)
	mux.Handle("/api", http.StripPrefix("/api", json))
	return mux
}

// HTMLServer returns new HTML server
func HTMLServer(userRepo app.UserRepo, itemRepo app.ItemRepo) http.Handler {
	server := Server{
		authMw: &htmlAuthMw{
			userRepo: userRepo,
		},
		userHandler: htmlUserHandler(userRepo),
		itemHandler: htmlItemHandler(itemRepo),
		router:      mux.NewRouter(),
	}
	server.routes(true)
	return &server
}

// JSONServer returns new JSON server
func JSONServer(userRepo app.UserRepo, itemRepo app.ItemRepo) http.Handler {
	server := Server{
		authMw: &jsonAuthMw{
			userRepo: userRepo,
		},
		userHandler: jsonUserHandler(userRepo),
		itemHandler: jsonItemHandler(itemRepo),
		router:      mux.NewRouter(),
	}
	server.routes(false)
	return &server
}

// Server represents an http server
type Server struct {
	authMw      AuthMw
	userHandler *UserHandler
	itemHandler *ItemHandler
	router      *mux.Router
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) routes(webMode bool) {
	if webMode {
		s.router.Handle("/", http.RedirectHandler("/signin", http.StatusFound))
		s.router.HandleFunc("/signin", s.userHandler.ShowSignin).Methods("GET")
	}

	s.router.HandleFunc("/signin", s.userHandler.ProcessSignin).Methods("POST")
	s.router.Handle("/items", ApplyFunc(s.itemHandler.Index,
		s.authMw.SetUser, s.authMw.RequireUser)).Methods("GET")
	s.router.Handle("/items", ApplyFunc(s.itemHandler.Create,
		s.authMw.SetUser, s.authMw.RequireUser)).Methods("POST")

	if webMode {
		s.router.Handle("/items/new", ApplyFunc(s.itemHandler.New,
			s.authMw.SetUser, s.authMw.RequireUser)).Methods("GET")
	}
}
