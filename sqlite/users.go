package sqlite

import (
	"database/sql"
	"strings"
	app "useritem"
)

// UserRepo is a Sqlite specific implementation of the user repository
type UserRepo struct {
	DB *sql.DB
}

// ByEmail will look for a user with the same email address
// return *app.User and an error
// if not found, return app.ErrNotFound
// if any SQL-specific error happens, pass the error through
//
// ByEmail is NOT case sensitive
func (repo *UserRepo) ByEmail(email string) (*app.User, error) {
	// prepare user
	user := app.User{
		Email: strings.ToLower(email),
	}

	// query row and get user
	row := repo.DB.QueryRow("select id, name, token from users where email=?", user.Email)
	err := row.Scan(&user.ID, &user.Name, &user.Token)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, app.ErrNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

// ByToken will look for a user with the same token
// return *app.User and an error
// if not found, return app.ErrNotFound
// if any SQL-specific error happens, pass the error through
func (repo *UserRepo) ByToken(token int) (*app.User, error) {
	// prepare user
	user := app.User{
		Token: token,
	}

	// query row and get user
	row := repo.DB.QueryRow("select id, name, email from users where token=?", user.Token)
	err := row.Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, app.ErrNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

// UpdateToken will update the token of a user with a specific id
// return an error
func (repo *UserRepo) UpdateToken(userID int, newToken int) error {
	_, err := repo.DB.Exec("update users set token=? where id=?", newToken, userID)
	return err
}
