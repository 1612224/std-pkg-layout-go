package http

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	app "useritem"
	"useritem/context"
)

type htmlAuthMw struct {
	userRepo app.UserRepo
}

// SetUser retrieves a user from session
// and put it into request context
func (a *htmlAuthMw) SetUser(next http.Handler) http.HandlerFunc {
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
func (a *htmlAuthMw) RequireUser(next http.Handler) http.HandlerFunc {
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

func htmlUserHandler(userRepo app.UserRepo) *UserHandler {
	uh := UserHandler{
		userRepo: userRepo,
		renderSignin: func(w http.ResponseWriter) {
			html := `
			<!DOCTYPE html>
			<html lang="en">
				<form action="/signin" method="POST">
					<label for="email">Email Address</label>
					<input type="email" id="email" name="email" placeholder="you@example.com">

					<label for="password">Password</label>
					<input type="password" id="password" name="password" placeholder="something-secret">

					<button type="submit">Sign in</button>
				</form>
			</html>`
			fmt.Fprint(w, html)
		},
		parseEmailAndPassword: func(r *http.Request) (email, password string) {
			email = r.PostFormValue("email")
			password = r.PostFormValue("password")
			return email, password
		},
		renderProcessSigninSuccess: func(w http.ResponseWriter, r *http.Request, token int) {
			cookie := http.Cookie{
				Name:  "session",
				Value: strconv.Itoa(token),
			}
			http.SetCookie(w, &cookie)
			http.Redirect(w, r, "/items", http.StatusFound)
		},
		renderProcessSigninError: func(w http.ResponseWriter, r *http.Request, err error) {
			switch err {
			case errAuthFailed:
				http.Redirect(w, r, "/signin", http.StatusFound)
			default:
				http.Error(w, "Something went wrong. Try again later.", http.StatusInternalServerError)
			}
		},
	}
	return &uh
}

func htmlItemHandler(itemRepo app.ItemRepo) *ItemHandler {
	ih := ItemHandler{
		itemRepo: itemRepo,
		renderNew: func(w http.ResponseWriter) {
			html := `
			<!DOCTYPE html>
			<html lang="en">
				<form action="/items" method="POST">
					<label for="name">Name</label>
					<input type="text" id="name" name="name" placeholder="Stop Item">

					<label for="price">Price</label>
					<input type="number" id="price" name="price" placeholder="18">

					<button type="submit">Create it!</button>
				</form>
			</html>`
			fmt.Fprint(w, html)
		},
		parseItem: func(r *http.Request) (*app.Item, error) {
			// Parse form values
			user := context.User(r.Context())
			item := app.Item{
				UserID: user.ID,
				Name:   r.PostFormValue("name"),
			}
			var err error
			item.Price, err = strconv.Atoi(r.PostFormValue("price"))
			return &item, err
		},
		renderCreateSuccess: func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/items", http.StatusFound)
		},
		renderCreateError: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, "Something went wrong. Try again later.", http.StatusInternalServerError)
		},
		renderIndexSuccess: func(w http.ResponseWriter, r *http.Request, items []app.Item) error {
			tplStr := `
			<!DOCTYPE html>
			<html lang="en">
				<h1>Items</h1>

				<ul>
				{{range .}}
				<li> {{.Name}}: <b>{{.Price}}VNĐ</b></li>
				{{end}}
				</ul>

				<p>
				<a href="/items/new">Create a new item</a>
				</p>
			</html>`
			tpl := template.Must(template.New("").Parse(tplStr))
			err := tpl.Execute(w, items)
			return err
		},
		renderIndexError: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		},
	}
	return &ih
}
