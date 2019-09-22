package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	app "useritem"
	"useritem/sqlite"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

var (
	db       *sql.DB
	userRepo app.UserRepo
	itemRepo app.ItemRepo
)

func main() {
	// setup db connection
	var err error
	db, err = sql.Open("sqlite3", "../database.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	// setup repos
	userRepo = &sqlite.UserRepo{DB: db}
	itemRepo = &sqlite.ItemRepo{DB: db}

	r := mux.NewRouter()
	r.Handle("/", http.RedirectHandler("/signin", http.StatusFound))
	r.HandleFunc("/signin", showSignin).Methods("GET")
	r.HandleFunc("/signin", processSignin).Methods("POST")
	r.HandleFunc("/items", allItems).Methods("GET")
	r.HandleFunc("/items", createItem).Methods("POST")
	r.HandleFunc("/items/new", newItem).Methods("GET")

	log.Println("Listening at port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func showSignin(w http.ResponseWriter, r *http.Request) {
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
}

func processSignin(w http.ResponseWriter, r *http.Request) {
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")

	// Lookup the user by their email in the DB
	email = strings.ToLower(email)
	row := db.QueryRow(`select id, password from users where email=?;`, email)
	var id int
	var truePassword string
	err := row.Scan(&id, &truePassword)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			// Email doesn't map to a user in our DB
			http.Redirect(w, r, "/signin", http.StatusFound)
		default:
			http.Error(w, "Something went wrong. Try again later.", http.StatusInternalServerError)
		}
		return
	}

	if truePassword != password {
		http.Redirect(w, r, "/signin", http.StatusNotFound)
		return
	}

	// Create a fake session token
	tokenStr := fmt.Sprintf("2019%d", id)
	token, _ := strconv.Atoi(tokenStr)
	err = userRepo.UpdateToken(id, token)
	if err != nil {
		http.Error(w, "Something went wrong. Try again later.", http.StatusInternalServerError)
		return
	}
	cookie := http.Cookie{
		Name:  "session",
		Value: tokenStr,
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/items", http.StatusFound)
}

func allItems(w http.ResponseWriter, r *http.Request) {
	// Verify the user is signed in
	session, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	token, _ := strconv.Atoi(session.Value)
	row := db.QueryRow(`select id from users where token=?;`, token)
	var userID int
	err = row.Scan(&userID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			// Email doesn't map to a user in our DB
			http.Redirect(w, r, "/signin", http.StatusFound)
		default:
			http.Error(w, "Something went wrong. Try again later.", http.StatusInternalServerError)
		}
		return
	}

	// Query for this user's items
	items, err := itemRepo.ByUser(userID)

	// Render the items
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
	err = tpl.Execute(w, items)
	if err != nil {
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
}

func createItem(w http.ResponseWriter, r *http.Request) {
	// Verify the user is signed in
	session, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	token, _ := strconv.Atoi(session.Value)
	row := db.QueryRow(`select id from users where token=?;`, token)
	var userID int
	err = row.Scan(&userID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			// Email doesn't map to a user in our DB
			http.Redirect(w, r, "/signin", http.StatusFound)
		default:
			http.Error(w, "Something went wrong. Try again later.", http.StatusInternalServerError)
		}
		return
	}

	// Parse form values and validate data
	name := r.PostFormValue("name")
	price, err := strconv.Atoi(r.PostFormValue("price"))
	if err != nil {
		http.Error(w, "Invalid price", http.StatusBadRequest)
		return
	}
	if price > 100000 {
		http.Error(w, "Price must be at 100,000 at maximum", http.StatusBadRequest)
	}

	// Create a new item
	item := app.Item{
		UserID: userID,
		Name:   name,
		Price:  price,
	}
	err = itemRepo.Create(&item)
	if err != nil {
		http.Error(w, "Something went wrong. Try again later.", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/items", http.StatusFound)
}

func newItem(w http.ResponseWriter, r *http.Request) {
	// Ignore auth for now - do it on the POST
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
}
