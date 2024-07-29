package usecases

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
)

const characters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type repository interface {
	domain.Transactor

	GetContender(ctx context.Context, tx domain.Transaction, contenderID domain.ResourceID) (domain.Contender, error)
	GetContenderByCode(ctx context.Context, tx domain.Transaction, registrationCode string) (domain.Contender, error)
	GetContendersByCompClass(ctx context.Context, tx domain.Transaction, compClassID domain.ResourceID) ([]domain.Contender, error)
	GetContendersByContest(ctx context.Context, tx domain.Transaction, contestID domain.ResourceID) ([]domain.Contender, error)
	StoreContender(ctx context.Context, tx domain.Transaction, contender domain.Contender) (domain.Contender, error)
	DeleteContender(ctx context.Context, tx domain.Transaction, contenderID domain.ResourceID) error
	GetContest(ctx context.Context, tx domain.Transaction, contestID domain.ResourceID) (domain.Contest, error)
	GetCompClass(ctx context.Context, tx domain.Transaction, compClassID domain.ResourceID) (domain.CompClass, error)
}

type ContenderUseCase struct {
	Repo        repository
	Authorizer  domain.Authorizer
	EventBroker domain.EventBroker
	ScoreKeeper domain.ScoreKeeper
}

func (uc *ContenderUseCase) GetContender(ctx context.Context, contenderID domain.ResourceID) (domain.Contender, error) {
	contender, err := uc.Repo.GetContender(ctx, nil, contenderID)
	if err != nil {
		return domain.Contender{}, err
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, contender.Ownership); err != nil {
		return domain.Contender{}, err
	}

	return withScore(contender, uc.ScoreKeeper), nil
}

func (uc *ContenderUseCase) GetContenderByCode(ctx context.Context, registrationCode string) (domain.Contender, error) {
	contender, err := uc.Repo.GetContenderByCode(ctx, nil, registrationCode)
	if err != nil {
		return domain.Contender{}, err
	}

	return withScore(contender, uc.ScoreKeeper), nil
}

func (uc *ContenderUseCase) GetContendersByCompClass(ctx context.Context, compClassID domain.ResourceID) ([]domain.Contender, error) {
	compClass, err := uc.Repo.GetCompClass(ctx, nil, compClassID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", domain.ErrNotFound, err)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, compClass.Ownership); err != nil {
		return nil, err
	}

	contenders, err := uc.Repo.GetContendersByCompClass(ctx, nil, compClassID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", domain.ErrRepositoryFailure, err)
	}

	for i, contender := range contenders {
		contenders[i] = withScore(contender, uc.ScoreKeeper)
	}

	return contenders, nil
}

func (uc *ContenderUseCase) GetContendersByContest(ctx context.Context, contestID domain.ResourceID) ([]domain.Contender, error) {
	contest, err := uc.Repo.GetContest(ctx, nil, contestID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", domain.ErrNotFound, err)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, contest.Ownership); err != nil {
		return nil, err
	}

	contenders, err := uc.Repo.GetContendersByContest(ctx, nil, contestID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", domain.ErrRepositoryFailure, err)
	}

	for i, contender := range contenders {
		contenders[i] = withScore(contender, uc.ScoreKeeper)
	}

	return contenders, nil
}

func (uc *ContenderUseCase) UpdateContender(ctx context.Context, contenderID domain.ResourceID, updates domain.Contender) (domain.Contender, error) {
	var mty domain.Contender
	var events []any

	contender, err := uc.Repo.GetContender(ctx, nil, contenderID)
	if err != nil {
		return mty, fmt.Errorf("%w: %w", domain.ErrNotFound, err)
	}

	role, err := uc.Authorizer.HasOwnership(ctx, contender.Ownership)
	if err != nil {
		return mty, err
	}

	publicInfoEvent := domain.ContenderPublicInfoUpdateEvent{
		ContenderID:         contenderID,
		CompClassID:         contender.CompClassID,
		PublicName:          contender.PublicName,
		ClubName:            contender.ClubName,
		WithdrawnFromFinals: contender.WithdrawnFromFinals,
		Disqualified:        contender.Disqualified,
	}

	publicInfoEventBaseline := publicInfoEvent

	contest, err := uc.Repo.GetContest(ctx, nil, contender.ContestID)
	if err != nil {
		return mty, fmt.Errorf("%w: %w", domain.ErrBadState, err)
	}

	if contender.CompClassID != updates.CompClassID && updates.CompClassID != 0 {
		var compClass domain.CompClass
		var err error

		publicInfoEvent.CompClassID = updates.CompClassID

		if contender.CompClassID != 0 {
			compClass, err = uc.Repo.GetCompClass(ctx, nil, contender.CompClassID)

			events = append(events, domain.ContenderSwitchClassEvent{
				ContenderID: contenderID,
				CompClassID: updates.CompClassID,
			})
		} else {
			compClass, err = uc.Repo.GetCompClass(ctx, nil, updates.CompClassID)

			events = append(events, domain.ContenderEnterEvent{
				ContenderID: contenderID,
				CompClassID: updates.CompClassID,
			})
		}

		if err != nil {
			return mty, fmt.Errorf("%w: %w", domain.ErrBadState, err)
		}

		gracePeriodEnd := compClass.TimeEnd.Add(contest.GracePeriod)

		if !role.OneOf(domain.AdminRole, domain.OrganizerRole) && time.Now().After(gracePeriodEnd) {
			return mty, domain.ErrContestEnded
		}
	}

	if contender.PublicName != updates.PublicName {
		publicInfoEvent.PublicName = updates.PublicName
	}

	if contender.ClubName != updates.ClubName {
		publicInfoEvent.ClubName = updates.ClubName
	}

	if contender.WithdrawnFromFinals != updates.WithdrawnFromFinals {
		publicInfoEvent.WithdrawnFromFinals = updates.WithdrawnFromFinals

		var event any
		if updates.WithdrawnFromFinals {
			event = domain.ContenderWithdrawFromFinalsEvent{
				ContenderID: contenderID,
			}
		} else {
			event = domain.ContenderReenterFinalsEvent{
				ContenderID: contenderID,
			}
		}

		events = append(events, event)
	}

	if contender.Disqualified != updates.Disqualified {
		publicInfoEvent.Disqualified = updates.Disqualified

		var event any

		if updates.Disqualified {
			event = domain.ContenderDisqualifyEvent{
				ContenderID: contenderID,
			}
		} else {
			event = domain.ContenderRequalifyEvent{
				ContenderID: contenderID,
			}
		}

		events = append(events, event)
	}

	if publicInfoEvent != publicInfoEventBaseline {
		events = append(events, publicInfoEvent)
	}

	contender.CompClassID = updates.CompClassID
	contender.Name = updates.Name
	contender.PublicName = updates.PublicName
	contender.ClubName = updates.ClubName
	contender.WithdrawnFromFinals = updates.WithdrawnFromFinals
	contender.Disqualified = updates.Disqualified

	if contender, err = uc.Repo.StoreContender(ctx, nil, contender); err != nil {
		return mty, fmt.Errorf("%w: %w", domain.ErrRepositoryFailure, err)
	}

	for _, event := range events {
		uc.EventBroker.Dispatch(contest.ID, event)
	}

	return withScore(contender, uc.ScoreKeeper), nil
}

func (uc *ContenderUseCase) DeleteContender(ctx context.Context, contenderID domain.ResourceID) error {
	contender, err := uc.Repo.GetContender(ctx, nil, contenderID)
	if err != nil {
		return domain.ErrNotFound
	}

	role, err := uc.Authorizer.HasOwnership(ctx, contender.Ownership)
	if err != nil {
		return err
	}

	if !role.OneOf(domain.AdminRole, domain.OrganizerRole) {
		return domain.ErrNotAllowed
	}

	if err := uc.Repo.DeleteContender(ctx, nil, contenderID); err != nil {
		return fmt.Errorf("%w: %w", domain.ErrRepositoryFailure, err)
	}

	return nil
}

const registrationCodeLength = 8

func (uc *ContenderUseCase) CreateContenders(ctx context.Context, contestID domain.ResourceID, number int) ([]domain.Contender, error) {
	contest, err := uc.Repo.GetContest(ctx, nil, contestID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", domain.ErrNotFound, err)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, contest.Ownership); err != nil {
		return nil, err
	}

	contenders := make([]domain.Contender, 0)

	tx := uc.Repo.Begin()
	defer tx.Rollback()

	for range number {
		var code []rune

		for range registrationCodeLength {
			code = append(code, []rune(characters)[rand.Intn(len(characters))])
		}

		contender := domain.Contender{
			ContestID: contestID,
			Ownership: domain.OwnershipData{
				OrganizerID: contest.Ownership.OrganizerID,
			},
			RegistrationCode: string(code),
		}

		if contender, err = uc.Repo.StoreContender(ctx, tx, contender); err != nil {
			return nil, fmt.Errorf("%w: %w", domain.ErrRepositoryFailure, err)
		}

		contenders = append(contenders, contender)
	}

	tx.Commit()

	return contenders, err
}
