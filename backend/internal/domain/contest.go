package domain

import (
	"context"
)

type ContestUseCase interface {
	GetContest(ctx context.Context, contestID ContestID) (Contest, error)
	GetContestsByOrganizer(ctx context.Context, organizerID OrganizerID) ([]Contest, error)
	UpdateContest(ctx context.Context, contestID ContestID, contest Contest) (Contest, error)
	DeleteContest(ctx context.Context, contestID ContestID) error
	DuplicateContest(ctx context.Context, contestID ContestID) (Contest, error)
	CreateContest(ctx context.Context, organizerID OrganizerID, contest Contest) (Contest, error)
	GetScoreboard(ctx context.Context, contestID ContestID) ([]ScoreboardEntry, error)
}
