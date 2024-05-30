package domain

type User struct {
	ID       ResourceID
	Name     string
	Username string
	Admin    bool
}

type UserUsecase interface {
}
