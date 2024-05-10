package models

type User struct {
	ID       int
	Username string
}

type Balance struct {
	ID      int
	UserID  int
	Balance int
}
