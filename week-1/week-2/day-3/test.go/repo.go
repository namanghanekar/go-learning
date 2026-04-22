package main

// Interface (important for mocking)
type UserRepo interface {
	GetUser(id int) string
}

// Real database (example)
type RealRepo struct{}

func (r RealRepo) GetUser(id int) string {
	return "Real User"
}
