package http

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	app "useritem"
)

var (
	errAuthFailed = errors.New("http: authentication failed")
)

// UserHandler handles an user session
type UserHandler struct {
	userRepo app.UserRepo

	renderSignin func(http.ResponseWriter)

	parseEmailAndPassword      func(*http.Request) (email, password string)
	renderProcessSigninSuccess func(w http.ResponseWriter, r *http.Request, token int)
	renderProcessSigninError   func(http.ResponseWriter, *http.Request, error)
}

// ShowSignin return signin page
func (h *UserHandler) ShowSignin(w http.ResponseWriter, r *http.Request) {
	h.renderSignin(w)
}

// ProcessSignin check signin credentials
func (h *UserHandler) ProcessSignin(w http.ResponseWriter, r *http.Request) {
	// Parse email & password
	email, password := h.parseEmailAndPassword(r)
	// Lookup the user by their email in the DB
	user, err := h.userRepo.ByEmail(email)
	if err != nil {
		switch err {
		case app.ErrNotFound:
			// Email doesn't map to a user in our DB
			h.renderProcessSigninError(w, r, errAuthFailed)
		default:
			log.Println(err)
			h.renderProcessSigninError(w, r, err)
		}
		return
	}

	// Check password
	if !user.CheckPassword(password) {
		http.Redirect(w, r, "/signin", http.StatusNotFound)
		return
	}

	// Create a fake session token
	tokenStr := fmt.Sprintf("2019%d", user.ID)
	token, _ := strconv.Atoi(tokenStr)
	err = h.userRepo.UpdateToken(user.ID, token)
	if err != nil {
		log.Println(err)
		h.renderProcessSigninError(w, r, err)
		return
	}
	h.renderProcessSigninSuccess(w, r, token)
}
