package http

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	app "useritem"
	"useritem/context"
)

// ItemHandler handles item related stuffs
type ItemHandler struct {
	itemRepo app.ItemRepo
}

// AllItems shows all items of an user
func (h *ItemHandler) AllItems(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())

	// Query for this user's items
	items, err := h.itemRepo.ByUser(user.ID)

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
		log.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
}

// CreateItem puts new item into item repo
func (h *ItemHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())

	// Parse form values and validate data
	item := app.Item{
		UserID: user.ID,
		Name:   r.PostFormValue("name"),
	}
	var err error
	item.Price, err = strconv.Atoi(r.PostFormValue("price"))
	if err != nil {
		http.Error(w, "Invalid price", http.StatusBadRequest)
		return
	}
	if item.Price > 100000 {
		http.Error(w, "Price must be at 100,000 at maximum", http.StatusBadRequest)
	}

	// Create a new item
	err = h.itemRepo.Create(&item)
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong. Try again later.", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/items", http.StatusFound)
}

// NewItem shows create new item page
func (h *ItemHandler) NewItem(w http.ResponseWriter, r *http.Request) {
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
