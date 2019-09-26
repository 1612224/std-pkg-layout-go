package app

// User represent an user's information
type User struct {
	ID       int
	Name     string
	Email    string
	Token    int
	password string
}

// SetPassword sets user password
func (u *User) SetPassword(password string) {
	u.password = password
}

// CheckPassword checks if a password is user's password
func (u *User) CheckPassword(password string) bool {
	if u.password != password {
		return false
	}
	return true
}

// Item is something that an user possesses
type Item struct {
	UserID int
	Name   string
	Price  int
}
