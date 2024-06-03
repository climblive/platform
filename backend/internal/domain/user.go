package domain

import "context"

type User struct {
	ID       ResourceID
	Name     string
	Username string
	Admin    bool
}

type UserUsecase interface {
	GetUser(ctx context.Context, id ResourceID) (User, error)
	GetCurrentUser(ctx context.Context) (User, error)
	GetUsers(ctx context.Context, organizerID ResourceID) ([]User, error)
	UpdateUser(ctx context.Context, id ResourceID, user User) (User, error)
	DeleteUser(ctx context.Context, id ResourceID) error
}
