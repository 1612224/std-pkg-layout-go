package http

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	app "useritem"
)

// UserHandler handles an user session
type UserHandler struct {
	userRepo app.UserRepo
}

// ShowSignin return signin page
func (h *UserHandler) ShowSignin(w http.ResponseWriter, r *http.Request) {
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

// ProcessSignin check signin credentials
func (h *UserHandler) ProcessSignin(w http.ResponseWriter, r *http.Request) {
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")

	// Lookup the user by their email in the DB
	user, err := h.userRepo.ByEmail(email)
	if err != nil {
		switch err {
		case app.ErrNotFound:
			// Email doesn't map to a user in our DB
			http.Redirect(w, r, "/signin", http.StatusFound)
		default:
			log.Println(err)
			http.Error(w, "Something went wrong. Try again later.", http.StatusInternalServerError)
		}
		return
	}

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
