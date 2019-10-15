package http

import (
	"log"
	"net/http"
	app "useritem"
	"useritem/context"
)

// ItemHandler handles item related stuffs
type ItemHandler struct {
	itemRepo app.ItemRepo

	renderNew func(http.ResponseWriter)

	parseItem           func(*http.Request) (*app.Item, error)
	renderCreateSuccess func(http.ResponseWriter, *http.Request)
	renderCreateError   func(http.ResponseWriter, *http.Request, error)

	renderIndexSuccess func(http.ResponseWriter, *http.Request, []app.Item) error
	renderIndexError   func(http.ResponseWriter, *http.Request, error)
}

// Index shows all items of an user
func (h *ItemHandler) Index(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())

	// Query for this user's items
	items, err := h.itemRepo.ByUser(user.ID)

	// Render the items
	if err != nil {
		h.renderIndexError(w, r, err)
		return
	}

	err = h.renderIndexSuccess(w, r, items)
	if err != nil {
		log.Println(err)
		h.renderIndexError(w, r, err)
	}
}

// Create puts new item into item repo
func (h *ItemHandler) Create(w http.ResponseWriter, r *http.Request) {

	// Parse item and validate data
	item, err := h.parseItem(r)
	if err != nil {
		h.renderCreateError(w, r, err)
		return
	}
	if item.Price > 100000 {
		h.renderCreateError(w, r, validationError{
			fields:  []string{"price"},
			message: "Price must be at most 100,000",
		})
		return
	}

	// Push new item into repo
	err = h.itemRepo.Create(item)
	if err != nil {
		log.Println(err)
		h.renderCreateError(w, r, err)
		return
	}
	h.renderCreateSuccess(w, r)
}

// New shows create new item page
func (h *ItemHandler) New(w http.ResponseWriter, r *http.Request) {
	// Ignore auth for now - do it on the POST
	h.renderNew(w)
}
