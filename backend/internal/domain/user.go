package domain

import "context"

type User struct {
	ID         ResourceID
	Name       string
	Username   string
	Admin      bool
	Organizers []ResourceID
}

type UserUseCase interface {
	GetCurrentUser(ctx context.Context) (User, error)
	GetUsersByOrganizer(ctx context.Context, organizerID ResourceID) ([]User, error)
	UpdateUser(ctx context.Context, userID ResourceID, user User) (User, error)
	DeleteUser(ctx context.Context, userID ResourceID) error
}
