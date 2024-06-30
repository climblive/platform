package usecases

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
)

type repository interface {
	GetContender(ctx context.Context, contenderID domain.ResourceID) (domain.Contender, error)
	GetContenderByCode(ctx context.Context, registrationCode string) (domain.Contender, error)
	GetContendersByCompClass(ctx context.Context, compClassID domain.ResourceID) ([]domain.Contender, error)
	GetContendersByContest(ctx context.Context, contestID domain.ResourceID) ([]domain.Contender, error)
	StoreContender(ctx context.Context, contender domain.Contender) error
	DeleteContender(ctx context.Context, contenderID domain.ResourceID) error
	GetContest(ctx context.Context, contestID domain.ResourceID) (domain.Contest, error)
	GetCompClass(ctx context.Context, compClassID domain.ResourceID) (domain.CompClass, error)
}

type ContenderUseCase struct {
	repo        repository
	authorizer  domain.Authorizer
	eventBroker domain.EventBroker
}

func NewContenderUseCase() domain.ContenderUsecase {
	return &ContenderUseCase{}
}

func (uc *ContenderUseCase) GetContender(ctx context.Context, contenderID domain.ResourceID) (domain.Contender, error) {
	contender, err := uc.repo.GetContender(ctx, contenderID)
	if err != nil {
		return domain.Contender{}, fmt.Errorf("%w: %w", domain.ErrNotFound, err)
	}

	if _, err := uc.authorizer.HasPermission(ctx, contender.Ownership); err != nil {
		return domain.Contender{}, err
	}

	return contender, nil
}

func (uc *ContenderUseCase) GetContenderByCode(ctx context.Context, registrationCode string) (domain.Contender, error) {
	return uc.repo.GetContenderByCode(ctx, registrationCode)
}

func (uc *ContenderUseCase) GetContendersByCompClass(ctx context.Context, compClassID domain.ResourceID) ([]domain.Contender, error) {
	compClass, err := uc.repo.GetCompClass(ctx, compClassID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", domain.ErrNotFound, err)
	}

	if _, err := uc.authorizer.HasPermission(ctx, compClass.Ownership); err != nil {
		return nil, err
	}

	contenders, err := uc.repo.GetContendersByCompClass(ctx, compClassID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", domain.ErrRepositoryFailure, err)
	}

	return contenders, nil
}

func (uc *ContenderUseCase) GetContendersByContest(ctx context.Context, contestID domain.ResourceID) ([]domain.Contender, error) {
	contest, err := uc.repo.GetContest(ctx, contestID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", domain.ErrNotFound, err)
	}

	if _, err := uc.authorizer.HasPermission(ctx, contest.Ownership); err != nil {
		return nil, err
	}

	contenders, err := uc.repo.GetContendersByContest(ctx, contestID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", domain.ErrRepositoryFailure, err)
	}

	return contenders, nil
}

func (uc *ContenderUseCase) UpdateContender(ctx context.Context, contenderID domain.ResourceID, upd domain.Contender) (domain.Contender, error) {
	var mty domain.Contender
	var events []any

	orig, err := uc.repo.GetContender(ctx, contenderID)
	if err != nil {
		return mty, fmt.Errorf("%w: %w", domain.ErrNotFound, err)
	}

	role, err := uc.authorizer.HasPermission(ctx, orig.Ownership)
	if err != nil {
		return mty, err
	}

	publicInfoEvent := domain.ContenderPublicInfoUpdateEvent{
		ContenderID:         contenderID,
		CompClassID:         orig.CompClassID,
		PublicName:          orig.PublicName,
		ClubName:            orig.ClubName,
		WithdrawnFromFinals: orig.WithdrawnFromFinals,
		Disqualified:        orig.Disqualified,
	}

	publicInfoEventBaseline := publicInfoEvent

	contest, err := uc.repo.GetContest(ctx, orig.ContestID)
	if err != nil {
		return mty, fmt.Errorf("%w: %w", domain.ErrBadState, err)
	}

	if orig.CompClassID != upd.CompClassID && upd.CompClassID != 0 {
		var compClass domain.CompClass
		var err error

		publicInfoEvent.CompClassID = upd.CompClassID

		if orig.CompClassID != 0 {
			compClass, err = uc.repo.GetCompClass(ctx, orig.CompClassID)

			events = append(events, domain.ContenderSwitchClassEvent{
				ContenderID: contenderID,
				CompClassID: upd.CompClassID,
			})
		} else {
			compClass, err = uc.repo.GetCompClass(ctx, upd.CompClassID)

			events = append(events, domain.ContenderEnterEvent{
				ContenderID: contenderID,
				CompClassID: upd.CompClassID,
			})
		}

		if err != nil {
			return mty, fmt.Errorf("%w: %w", domain.ErrBadState, err)
		}

		gracePeriodEnd := compClass.TimeEnd.Add(contest.GracePeriod)

		if role != domain.OrganizerRole && time.Now().After(gracePeriodEnd) {
			return mty, domain.ErrContestEnded
		}
	}

	if orig.PublicName != upd.PublicName {
		publicInfoEvent.PublicName = upd.PublicName
	}

	if orig.ClubName != upd.ClubName {
		publicInfoEvent.ClubName = upd.ClubName
	}

	if orig.WithdrawnFromFinals != upd.WithdrawnFromFinals {
		publicInfoEvent.WithdrawnFromFinals = upd.WithdrawnFromFinals

		var event any
		if upd.WithdrawnFromFinals {
			event = domain.ContenderWithdrawFromFinalsEvent{
				ContenderID: contenderID,
			}
		} else {
			event = domain.ContenderReenterIntoFinalsEvent{
				ContenderID: contenderID,
			}
		}

		events = append(events, event)
	}

	if orig.Disqualified != upd.Disqualified {
		publicInfoEvent.Disqualified = upd.Disqualified

		var event any

		if upd.Disqualified {
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

	orig.CompClassID = upd.CompClassID
	orig.Name = upd.Name
	orig.PublicName = upd.PublicName
	orig.ClubName = upd.ClubName
	orig.WithdrawnFromFinals = upd.WithdrawnFromFinals
	orig.Disqualified = upd.Disqualified

	if err := uc.repo.StoreContender(ctx, orig); err != nil {
		return mty, fmt.Errorf("%w: %w", domain.ErrRepositoryFailure, err)
	}

	for _, event := range events {
		uc.eventBroker.Dispatch(contest.ID, event)
	}

	return orig, nil
}

func (uc *ContenderUseCase) DeleteContender(ctx context.Context, contenderID domain.ResourceID) error {
	contender, err := uc.repo.GetContender(ctx, contenderID)
	if err != nil {
		return domain.ErrNotFound
	}

	role, err := uc.authorizer.HasPermission(ctx, contender.Ownership)
	if err != nil {
		return nil
	}

	if !role.OneOf(domain.AdminRole, domain.OrganizerRole) {
		return domain.ErrNotAllowed
	}

	if err := uc.repo.DeleteContender(ctx, contenderID); err != nil {
		return fmt.Errorf("%w: %w", domain.ErrRepositoryFailure, err)
	}

	return nil
}

const registrationCodeLength = 8

func (uc *ContenderUseCase) CreateContenders(ctx context.Context, contestID domain.ResourceID, number int) ([]domain.Contender, error) {
	contest, err := uc.repo.GetContest(ctx, contestID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", domain.ErrNotFound, err)
	}

	if _, err := uc.authorizer.HasPermission(ctx, contest.Ownership); err != nil {
		return nil, err
	}

	// TODO: begin transaction
	characters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

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

		if err := uc.repo.StoreContender(ctx, contender); err != nil {
			return []domain.Contender{}, fmt.Errorf("%w: %w", domain.ErrRepositoryFailure, err)
		}
	}

	return nil, nil
}
