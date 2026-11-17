package usecases_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/usecases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateCompetitionAttemptPriority(t *testing.T) {
	for _, favorTopAttempts := range []bool{false, true} {
		t.Run(fmt.Sprintf("favorTopAttempts=%t", favorTopAttempts), func(t *testing.T) {
			repo := new(repositoryMock)
			authorizer := new(authorizerMock)
			ownership := domain.OwnershipData{OrganizerID: 1}
			repo.On("GetOrganizer", mock.Anything, nil, domain.OrganizerID(1)).
				Return(domain.Organizer{ID: 1, Ownership: ownership}, nil)
			authorizer.On("HasOwnership", mock.Anything, ownership).Return(domain.OrganizerRole, nil)
			repo.On("GetCompetitionsByOrganizer", mock.Anything, nil, domain.OrganizerID(1)).
				Return([]domain.Competition{}, nil)
			repo.On("StoreCompetition", mock.Anything, nil, mock.MatchedBy(func(competition domain.Competition) bool {
				return competition.FavorTopAttempts == favorTopAttempts && !competition.UsePoints
			})).Return(domain.Competition{ID: 2, FavorTopAttempts: favorTopAttempts}, nil)
			uc := usecases.CompetitionUseCase{Repo: repo, Authorizer: authorizer}
			competition, err := uc.CreateCompetition(context.Background(), 1, domain.CompetitionTemplate{
				Name: "Test", Country: "SE", NameRetentionTime: 14 * 24 * time.Hour,
				FavorTopAttempts: favorTopAttempts,
			})
			require.NoError(t, err)
			assert.Equal(t, favorTopAttempts, competition.FavorTopAttempts)
			repo.AssertExpectations(t)
			authorizer.AssertExpectations(t)
		})
	}
}

func TestPatchCompetitionAttemptPriority(t *testing.T) {
	tests := []struct {
		name    string
		initial bool
		patch   domain.Patch[bool]
		updated bool
	}{
		{name: "favor attempts", patch: domain.NewPatch(true), updated: true},
		{name: "favor zones", initial: true, patch: domain.NewPatch(false)},
		{name: "unchanged", initial: true, patch: domain.NewPatch(true), updated: true},
		{name: "omitted preserves legacy", initial: true, updated: true},
		{name: "omitted preserves zones"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(repositoryMock)
			authorizer := new(authorizerMock)
			broker := new(eventBrokerMock)
			competition := domain.Competition{
				ID: 1, Ownership: domain.OwnershipData{OrganizerID: 1}, Name: "Test", Country: "SE",
				NameRetentionTime: 14 * 24 * time.Hour, FavorTopAttempts: tt.initial,
			}
			repo.On("GetCompetition", mock.Anything, nil, competition.ID).Return(competition, nil)
			authorizer.On("HasOwnership", mock.Anything, competition.Ownership).Return(domain.OrganizerRole, nil)
			expected := competition
			expected.FavorTopAttempts = tt.updated
			repo.On("StoreCompetition", mock.Anything, nil, expected).Return(expected, nil)
			if tt.initial != tt.updated {
				broker.On("Dispatch", competition.ID, domain.RulesUpdatedEvent{
					CompetitionID: competition.ID, FavorTopAttempts: tt.updated,
				}).Return(nil).Once()
			}
			uc := usecases.CompetitionUseCase{Repo: repo, Authorizer: authorizer, EventBroker: broker}
			updated, err := uc.PatchCompetition(context.Background(), competition.ID,
				domain.CompetitionPatch{FavorTopAttempts: tt.patch})
			require.NoError(t, err)
			assert.Equal(t, expected, updated)
			repo.AssertExpectations(t)
			authorizer.AssertExpectations(t)
			broker.AssertExpectations(t)
			if tt.initial == tt.updated {
				broker.AssertNotCalled(t, "Dispatch", mock.Anything, mock.Anything)
			}
		})
	}
}
