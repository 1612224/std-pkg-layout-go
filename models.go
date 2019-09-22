package app

// User represent an user's information
type User struct {
	ID    int
	Name  string
	Email string
	Token int
}

// Item is something that an user possesses
type Item struct {
	UserID int
	Name   string
	Price  int
}
