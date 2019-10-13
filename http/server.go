package http

import (
	"net/http"
	app "useritem"

	"github.com/gorilla/mux"
)

// NewServer returns new server
func NewServer(userRepo app.UserRepo, itemRepo app.ItemRepo) *Server {
	server := Server{
		authMw: &AuthMw{
			userRepo: userRepo,
		},
		userHandler: &UserHandler{
			userRepo: userRepo,
		},
		itemHandler: &ItemHandler{
			itemRepo: itemRepo,
		},
		router: mux.NewRouter(),
	}
	server.routes()
	return &server
}

// Server represents an http server
type Server struct {
	authMw      *AuthMw
	userHandler *UserHandler
	itemHandler *ItemHandler
	router      *mux.Router
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) routes() {
	s.router.Handle("/", http.RedirectHandler("/signin", http.StatusFound))
	s.router.HandleFunc("/signin", s.userHandler.ShowSignin).Methods("GET")
	s.router.HandleFunc("/signin", s.userHandler.ProcessSignin).Methods("POST")
	s.router.Handle("/items", ApplyFunc(s.itemHandler.AllItems,
		s.authMw.UserViaSession, s.authMw.RequireUser)).Methods("GET")
	s.router.Handle("/items", ApplyFunc(s.itemHandler.CreateItem,
		s.authMw.UserViaSession, s.authMw.RequireUser)).Methods("POST")
	s.router.Handle("/items/new", ApplyFunc(s.itemHandler.NewItem,
		s.authMw.UserViaSession, s.authMw.RequireUser)).Methods("GET")
}
