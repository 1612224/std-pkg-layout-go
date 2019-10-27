package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	app "useritem"
	"useritem/context"

	"golang.org/x/oauth2"
)

func renderJSON(w http.ResponseWriter, data interface{}, status int) {
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.Encode(data)
}

type jsonError struct {
	Message string `json:"error"`
	Type    string `json:"type"`
}

func (e jsonError) Error() string {
	return fmt.Sprintf("json %s error: %s", e.Type, e.Message)
}

type jsonAuthMw struct {
	userRepo app.UserRepo
}

// SetUser retrieves a user from session
// and put it into request context
func (mw *jsonAuthMw) SetUser(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bearer := r.Header.Get("Authorization")
		if len(bearer) <= len("Bearer") {
			next.ServeHTTP(w, r)
			return
		}
		tokenStr := strings.TrimSpace(bearer[len("Bearer"):])
		token, err := strconv.Atoi(tokenStr)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		user, err := mw.userRepo.ByToken(token)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		r = r.WithContext(context.WithUser(r.Context(), user))
		next.ServeHTTP(w, r)
	}
}

// RequireUser requires a user from context
// if no user found, redirect to sign in
func (mw *jsonAuthMw) RequireUser(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())
		if user == nil {
			renderJSON(w, jsonError{
				Message: "Unauthorized access. Do you have a valid oath2 token set ?",
				Type:    "unauthorized",
			}, http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func jsonUserHandler(userRepo app.UserRepo) *UserHandler {
	uh := UserHandler{
		userRepo: userRepo,

		parseEmailAndPassword: func(r *http.Request) (email, password string) {
			var req struct {
				Email    string `json:"email"`
				Password string `json:"password"`
			}
			dec := json.NewDecoder(r.Body)
			dec.Decode(&req)
			return req.Email, req.Password
		},
		renderProcessSigninSuccess: func(w http.ResponseWriter, r *http.Request, token int) {
			t := oauth2.Token{
				TokenType:   "Bearer",
				AccessToken: strconv.Itoa(token),
			}
			renderJSON(w, t, http.StatusOK)
		},
		renderProcessSigninError: func(w http.ResponseWriter, r *http.Request, err error) {
			switch err {
			case errAuthFailed:
				renderJSON(w, jsonError{
					Message: "Invalid authentication details",
					Type:    "authentication",
				}, http.StatusBadRequest)
			default:
				renderJSON(w, jsonError{
					Message: "Something went wrong. Try again later",
					Type:    "internal_server",
				}, http.StatusInternalServerError)
			}
		},
	}
	return &uh
}

type jsonItem struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
}

func (item *jsonItem) read(i app.Item) {
	item.Name = i.Name
	item.Price = i.Price
}

func jsonItemHandler(itemRepo app.ItemRepo) *ItemHandler {
	ih := ItemHandler{
		itemRepo: itemRepo,

		parseItem: func(r *http.Request) (*app.Item, error) {
			var req struct {
				Name  string `json:"name"`
				Price int    `json:"price"`
			}
			dec := json.NewDecoder(r.Body)
			err := dec.Decode(&req)
			if err != nil {
				return nil, validationError{
					fields:  []string{"price"},
					message: "Price must be integer",
				}
			}

			return &app.Item{
				Name:  req.Name,
				Price: req.Price,
			}, nil
		},
		renderCreateSuccess: func(w http.ResponseWriter, r *http.Request, item *app.Item) {
			var res jsonItem
			res.read(*item)
			renderJSON(w, res, http.StatusCreated)
		},
		renderCreateError: func(w http.ResponseWriter, r *http.Request, err error) {
			switch v := err.(type) {
			case validationError:
				renderJSON(w, struct {
					Fields []string `json:"fields"`
					jsonError
				}{
					Fields: v.fields,
					jsonError: jsonError{
						Message: v.message,
						Type:    "validation",
					},
				}, http.StatusBadRequest)
			default:
				renderJSON(w, jsonError{
					Message: "Something went wrong. Try again later",
					Type:    "internal_server",
				}, http.StatusInternalServerError)
			}
		},
		renderIndexSuccess: func(w http.ResponseWriter, r *http.Request, items []app.Item) error {
			res := make([]jsonItem, 0, len(items))
			for _, item := range items {
				var ji jsonItem
				ji.read(item)
				res = append(res, ji)
			}
			enc := json.NewEncoder(w)
			return enc.Encode(res)
		},
		renderIndexError: func(w http.ResponseWriter, r *http.Request, err error) {
			renderJSON(w, jsonError{
				Message: "Something went wrong. Try again later",
				Type:    "internal_server",
			}, http.StatusInternalServerError)
		},
	}
	return &ih
}
