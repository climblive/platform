package usecases

import (
	"context"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

type unlockRequestUseCaseRepository interface {
	domain.Transactor

	GetContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) (domain.Contest, error)
	GetOrganizer(ctx context.Context, tx domain.Transaction, organizerID domain.OrganizerID) (domain.Organizer, error)
	CreateUnlockRequest(ctx context.Context, tx domain.Transaction, template domain.UnlockRequestTemplate, organizerID domain.OrganizerID) (domain.UnlockRequestID, error)
	GetUnlockRequest(ctx context.Context, tx domain.Transaction, id domain.UnlockRequestID) (domain.UnlockRequest, error)
	GetUnlockRequestsByContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) ([]domain.UnlockRequest, error)
	GetUnlockRequestsByOrganizer(ctx context.Context, tx domain.Transaction, organizerID domain.OrganizerID) ([]domain.UnlockRequest, error)
	GetPendingUnlockRequests(ctx context.Context, tx domain.Transaction) ([]domain.UnlockRequest, error)
	UpdateUnlockRequestStatus(ctx context.Context, tx domain.Transaction, id domain.UnlockRequestID, review domain.UnlockRequestReview) error
	HasApprovedUnlockRequest(ctx context.Context, tx domain.Transaction, organizerID domain.OrganizerID) (bool, error)
	StoreContest(ctx context.Context, tx domain.Transaction, contest domain.Contest) (domain.Contest, error)
}

type UnlockRequestUseCase struct {
	Authorizer domain.Authorizer
	Repo       unlockRequestUseCaseRepository
}

func (uc *UnlockRequestUseCase) CreateUnlockRequest(ctx context.Context, template domain.UnlockRequestTemplate) (domain.UnlockRequestID, error) {
	contest, err := uc.Repo.GetContest(ctx, nil, template.ContestID)
	if err != nil {
		return 0, errors.Wrap(err, 0)
	}

	authRole, err := uc.Authorizer.HasOwnership(ctx, contest.Ownership)
	if err != nil {
		return 0, errors.Wrap(err, 0)
	}

	if !contest.EvaluationMode {
		return 0, errors.Wrap(domain.ErrNotAllowed, 0)
	}

	tx, err := uc.Repo.Begin()
	if err != nil {
		return 0, errors.Wrap(err, 0)
	}
	defer tx.Rollback()

	hasApprovedRequest, err := uc.Repo.HasApprovedUnlockRequest(ctx, tx, contest.Ownership.OrganizerID)
	if err != nil {
		return 0, errors.Wrap(err, 0)
	}

	requestID, err := uc.Repo.CreateUnlockRequest(ctx, tx, template, contest.Ownership.OrganizerID)
	if err != nil {
		return 0, errors.Wrap(err, 0)
	}

	if hasApprovedRequest && authRole == domain.OrganizerRole {
		review := domain.UnlockRequestReview{
			Status: domain.UnlockRequestStatusApproved,
		}
		if err := uc.Repo.UpdateUnlockRequestStatus(ctx, tx, requestID, review); err != nil {
			return 0, errors.Wrap(err, 0)
		}

		contest.EvaluationMode = false
		if _, err := uc.Repo.StoreContest(ctx, tx, contest); err != nil {
			return 0, errors.Wrap(err, 0)
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, errors.Wrap(err, 0)
	}

	return requestID, nil
}

func (uc *UnlockRequestUseCase) GetUnlockRequest(ctx context.Context, id domain.UnlockRequestID) (domain.UnlockRequest, error) {
	request, err := uc.Repo.GetUnlockRequest(ctx, nil, id)
	if err != nil {
		return domain.UnlockRequest{}, errors.Wrap(err, 0)
	}

	organizer, err := uc.Repo.GetOrganizer(ctx, nil, request.OrganizerID)
	if err != nil {
		return domain.UnlockRequest{}, errors.Wrap(err, 0)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, organizer.Ownership); err != nil {
		authRole, authErr := uc.Authorizer.HasOwnership(ctx, domain.OwnershipData{})
		if authErr != nil || authRole != domain.AdminRole {
			return domain.UnlockRequest{}, errors.Wrap(domain.ErrNotAuthorized, 0)
		}
	}

	return request, nil
}

func (uc *UnlockRequestUseCase) GetUnlockRequestsByContest(ctx context.Context, contestID domain.ContestID) ([]domain.UnlockRequest, error) {
	contest, err := uc.Repo.GetContest(ctx, nil, contestID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, contest.Ownership); err != nil {
		return nil, errors.Wrap(err, 0)
	}

	requests, err := uc.Repo.GetUnlockRequestsByContest(ctx, nil, contestID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	return requests, nil
}

func (uc *UnlockRequestUseCase) GetUnlockRequestsByOrganizer(ctx context.Context, organizerID domain.OrganizerID) ([]domain.UnlockRequest, error) {
	organizer, err := uc.Repo.GetOrganizer(ctx, nil, organizerID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, organizer.Ownership); err != nil {
		return nil, errors.Wrap(err, 0)
	}

	requests, err := uc.Repo.GetUnlockRequestsByOrganizer(ctx, nil, organizerID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	return requests, nil
}

func (uc *UnlockRequestUseCase) GetPendingUnlockRequests(ctx context.Context) ([]domain.UnlockRequest, error) {
	authRole, err := uc.Authorizer.HasOwnership(ctx, domain.OwnershipData{})
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	if authRole != domain.AdminRole {
		return nil, domain.ErrNotAuthorized
	}

	requests, err := uc.Repo.GetPendingUnlockRequests(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	return requests, nil
}

func (uc *UnlockRequestUseCase) ReviewUnlockRequest(ctx context.Context, id domain.UnlockRequestID, review domain.UnlockRequestReview) error {
	authRole, err := uc.Authorizer.HasOwnership(ctx, domain.OwnershipData{})
	if err != nil {
		return errors.Wrap(err, 0)
	}

	if authRole != domain.AdminRole {
		return domain.ErrNotAuthorized
	}

	tx, err := uc.Repo.Begin()
	if err != nil {
		return errors.Wrap(err, 0)
	}
	defer tx.Rollback()

	request, err := uc.Repo.GetUnlockRequest(ctx, tx, id)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	if request.Status != domain.UnlockRequestStatusPending {
		return errors.Wrap(domain.ErrNotAllowed, 0)
	}

	if err := uc.Repo.UpdateUnlockRequestStatus(ctx, tx, id, review); err != nil {
		return errors.Wrap(err, 0)
	}

	if review.Status == domain.UnlockRequestStatusApproved {
		contest, err := uc.Repo.GetContest(ctx, tx, request.ContestID)
		if err != nil {
			return errors.Wrap(err, 0)
		}

		contest.EvaluationMode = false
		if _, err := uc.Repo.StoreContest(ctx, tx, contest); err != nil {
			return errors.Wrap(err, 0)
		}
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, 0)
	}

	return nil
}
