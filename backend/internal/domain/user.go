package domain

import "context"

type UserUseCase interface {
	GetCurrentUser(ctx context.Context) (User, error)
	GetUsersByOrganizer(ctx context.Context, organizerID OrganizerID) ([]User, error)
	UpdateUser(ctx context.Context, userID UserID, user User) (User, error)
	DeleteUser(ctx context.Context, userID UserID) error
}
