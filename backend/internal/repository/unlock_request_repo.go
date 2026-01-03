package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/climblive/platform/backend/internal/database"
	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

func (d *Database) CreateUnlockRequest(ctx context.Context, tx domain.Transaction, template domain.UnlockRequestTemplate, requestedByUserID domain.UserID, organizerID domain.OrganizerID) (domain.UnlockRequestID, error) {
	reason := sql.NullString{}
	if template.Reason != "" {
		reason = sql.NullString{String: template.Reason, Valid: true}
	}

	id, err := d.WithTx(tx).CreateUnlockRequest(ctx, database.CreateUnlockRequestParams{
		ContestID:         int32(template.ContestID),
		OrganizerID:       int32(organizerID),
		RequestedByUserID: int32(requestedByUserID),
		Reason:            reason,
	})
	if err != nil {
		return 0, errors.Wrap(err, 0)
	}

	return domain.UnlockRequestID(id), nil
}

func (d *Database) GetUnlockRequest(ctx context.Context, tx domain.Transaction, id domain.UnlockRequestID) (domain.UnlockRequest, error) {
	record, err := d.WithTx(tx).GetUnlockRequest(ctx, int32(id))
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return domain.UnlockRequest{}, errors.Wrap(domain.ErrNotFound, 0)
	case err != nil:
		return domain.UnlockRequest{}, errors.Wrap(err, 0)
	}

	return unlockRequestToDomain(record), nil
}

func (d *Database) GetUnlockRequestsByContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) ([]domain.UnlockRequest, error) {
	records, err := d.WithTx(tx).GetUnlockRequestsByContest(ctx, int32(contestID))
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	requests := make([]domain.UnlockRequest, 0, len(records))
	for _, record := range records {
		requests = append(requests, unlockRequestToDomain(record))
	}

	return requests, nil
}

func (d *Database) GetUnlockRequestsByOrganizer(ctx context.Context, tx domain.Transaction, organizerID domain.OrganizerID) ([]domain.UnlockRequest, error) {
	records, err := d.WithTx(tx).GetUnlockRequestsByOrganizer(ctx, int32(organizerID))
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	requests := make([]domain.UnlockRequest, 0, len(records))
	for _, record := range records {
		requests = append(requests, unlockRequestToDomain(record))
	}

	return requests, nil
}

func (d *Database) GetPendingUnlockRequests(ctx context.Context, tx domain.Transaction) ([]domain.UnlockRequest, error) {
	records, err := d.WithTx(tx).GetPendingUnlockRequests(ctx)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	requests := make([]domain.UnlockRequest, 0, len(records))
	for _, record := range records {
		requests = append(requests, unlockRequestToDomain(record))
	}

	return requests, nil
}

func (d *Database) UpdateUnlockRequestStatus(ctx context.Context, tx domain.Transaction, id domain.UnlockRequestID, review domain.UnlockRequestReview, reviewedByUserID domain.UserID) error {
	reviewedAt := sql.NullTime{Time: time.Now(), Valid: true}
	reviewerID := sql.NullInt32{Int32: int32(reviewedByUserID), Valid: true}
	reviewNote := sql.NullString{}
	if review.ReviewNote != "" {
		reviewNote = sql.NullString{String: review.ReviewNote, Valid: true}
	}

	err := d.WithTx(tx).UpdateUnlockRequestStatus(ctx, database.UpdateUnlockRequestStatusParams{
		Status:             database.UnlockRequestStatus(review.Status),
		ReviewedByUserID:   reviewerID,
		ReviewedAt:         reviewedAt,
		ReviewNote:         reviewNote,
		ID:                 int32(id),
	})
	if err != nil {
		return errors.Wrap(err, 0)
	}

	return nil
}

func (d *Database) HasApprovedUnlockRequest(ctx context.Context, tx domain.Transaction, organizerID domain.OrganizerID) (bool, error) {
	result, err := d.WithTx(tx).HasApprovedUnlockRequest(ctx, int32(organizerID))
	if err != nil {
		return false, errors.Wrap(err, 0)
	}

	return result, nil
}

func unlockRequestToDomain(record database.UnlockRequest) domain.UnlockRequest {
	request := domain.UnlockRequest{
		ID:                domain.UnlockRequestID(record.ID),
		ContestID:         domain.ContestID(record.ContestID),
		OrganizerID:       domain.OrganizerID(record.OrganizerID),
		RequestedByUserID: domain.UserID(record.RequestedByUserID),
		Status:            domain.UnlockRequestStatus(record.Status),
		CreatedAt:         record.CreatedAt,
	}

	if record.ReviewedByUserID.Valid {
		userID := domain.UserID(record.ReviewedByUserID.Int32)
		request.ReviewedByUserID = &userID
	}

	if record.ReviewedAt.Valid {
		request.ReviewedAt = &record.ReviewedAt.Time
	}

	if record.Reason.Valid {
		request.Reason = record.Reason.String
	}

	if record.ReviewNote.Valid {
		request.ReviewNote = record.ReviewNote.String
	}

	return request
}
