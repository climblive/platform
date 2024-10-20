package domain

import "context"

type User struct {
	ID         UserID
	Name       string
	Username   string
	Admin      bool
	Organizers []OrganizerID
}

type UserUseCase interface {
	GetCurrentUser(ctx context.Context) (User, error)
	GetUsersByOrganizer(ctx context.Context, organizerID OrganizerID) ([]User, error)
	UpdateUser(ctx context.Context, userID UserID, user User) (User, error)
	DeleteUser(ctx context.Context, userID UserID) error
}
