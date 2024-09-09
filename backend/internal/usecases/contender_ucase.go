package usecases

import (
	"context"
	"strings"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

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
	Repo                      repository
	Authorizer                domain.Authorizer
	EventBroker               domain.EventBroker
	ScoreKeeper               domain.ScoreKeeper
	RegistrationCodeGenerator domain.CodeGenerator
}

func (uc *ContenderUseCase) GetContender(ctx context.Context, contenderID domain.ResourceID) (domain.Contender, error) {
	contender, err := uc.Repo.GetContender(ctx, nil, contenderID)
	if err != nil {
		return domain.Contender{}, errors.Wrap(err, 0)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, contender.Ownership); err != nil {
		return domain.Contender{}, errors.Wrap(err, 0)
	}

	return withScore(contender, uc.ScoreKeeper), nil
}

func (uc *ContenderUseCase) GetContenderByCode(ctx context.Context, registrationCode string) (domain.Contender, error) {
	contender, err := uc.Repo.GetContenderByCode(ctx, nil, registrationCode)
	if err != nil {
		return domain.Contender{}, errors.Wrap(err, 0)
	}

	return withScore(contender, uc.ScoreKeeper), nil
}

func (uc *ContenderUseCase) GetContendersByCompClass(ctx context.Context, compClassID domain.ResourceID) ([]domain.Contender, error) {
	compClass, err := uc.Repo.GetCompClass(ctx, nil, compClassID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, compClass.Ownership); err != nil {
		return nil, errors.Wrap(err, 0)
	}

	contenders, err := uc.Repo.GetContendersByCompClass(ctx, nil, compClassID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	for i, contender := range contenders {
		contenders[i] = withScore(contender, uc.ScoreKeeper)
	}

	return contenders, nil
}

func (uc *ContenderUseCase) GetContendersByContest(ctx context.Context, contestID domain.ResourceID) ([]domain.Contender, error) {
	contest, err := uc.Repo.GetContest(ctx, nil, contestID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, contest.Ownership); err != nil {
		return nil, errors.Wrap(err, 0)
	}

	contenders, err := uc.Repo.GetContendersByContest(ctx, nil, contestID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
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
		return mty, errors.Wrap(err, 0)
	}

	role, err := uc.Authorizer.HasOwnership(ctx, contender.Ownership)
	if err != nil {
		return mty, errors.Wrap(err, 0)
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
		return mty, errors.Errorf("%w: %w", domain.ErrRepositoryIntegrityViolation, err)
	}

	if contender.CompClassID != 0 {
		compClass, err := uc.Repo.GetCompClass(ctx, nil, contender.CompClassID)
		if err != nil {
			return mty, errors.Errorf("%w: %w", domain.ErrRepositoryIntegrityViolation, err)
		}

		gracePeriodEnd := compClass.TimeEnd.Add(contest.GracePeriod)
		switch {
		case role.OneOf(domain.AdminRole, domain.OrganizerRole):
			break;
		case time.Now().After(gracePeriodEnd):
			return mty, errors.Wrap(domain.ErrContestEnded, 0)
		}
	}

	if contender.CompClassID != updates.CompClassID {
		if updates.CompClassID == 0 {
			return mty, errors.Wrap(domain.ErrNotAllowed, 0)
		}

		compClass, err := uc.Repo.GetCompClass(ctx, nil, updates.CompClassID)
		if err != nil {
			return mty, errors.Wrap(err, 0)
		}

		if contender.CompClassID == 0 {
			events = append(events, domain.ContenderEnterEvent{
				ContenderID: contenderID,
				CompClassID: updates.CompClassID,
			})
		} else {
			events = append(events, domain.ContenderSwitchClassEvent{
				ContenderID: contenderID,
				CompClassID: updates.CompClassID,
			})
		}

		gracePeriodEnd := compClass.TimeEnd.Add(contest.GracePeriod)

		switch {
		case role.OneOf(domain.AdminRole, domain.OrganizerRole):
			break;
		case time.Now().After(gracePeriodEnd):
			return mty, errors.Wrap(domain.ErrContestEnded, 0)
		}

		contender.CompClassID = compClass.ID
		publicInfoEvent.CompClassID = updates.CompClassID

		if contender.Entered == nil {
			timestamp := time.Now()
			contender.Entered = &timestamp
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
		if !role.OneOf(domain.AdminRole, domain.OrganizerRole) {
			return mty, errors.Wrap(domain.ErrInsufficientRole, 0)
		}

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
	contender.Name = strings.TrimSpace(updates.Name)
	contender.PublicName = strings.TrimSpace(updates.PublicName)
	contender.ClubName = strings.TrimSpace(updates.ClubName)
	contender.WithdrawnFromFinals = updates.WithdrawnFromFinals
	contender.Disqualified = updates.Disqualified

	if contender.Name == "" {
		return mty, errors.Errorf("%w: %w", domain.ErrInvalidData, domain.ErrEmptyName)
	}

	if contender, err = uc.Repo.StoreContender(ctx, nil, contender); err != nil {
		return mty, errors.Wrap(err, 0)
	}

	for _, event := range events {
		uc.EventBroker.Dispatch(contest.ID, event)
	}

	return withScore(contender, uc.ScoreKeeper), nil
}

func (uc *ContenderUseCase) DeleteContender(ctx context.Context, contenderID domain.ResourceID) error {
	contender, err := uc.Repo.GetContender(ctx, nil, contenderID)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	role, err := uc.Authorizer.HasOwnership(ctx, contender.Ownership)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	if !role.OneOf(domain.AdminRole, domain.OrganizerRole) {
		return errors.Wrap(domain.ErrInsufficientRole, 0)
	}

	if err := uc.Repo.DeleteContender(ctx, nil, contenderID); err != nil {
		return errors.Wrap(err, 0)
	}

	return nil
}

const registrationCodeLength = 8

func (uc *ContenderUseCase) CreateContenders(ctx context.Context, contestID domain.ResourceID, number int) ([]domain.Contender, error) {
	contest, err := uc.Repo.GetContest(ctx, nil, contestID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, contest.Ownership); err != nil {
		return nil, errors.Wrap(err, 0)
	}

	contenders := make([]domain.Contender, 0)

	tx := uc.Repo.Begin()

	for range number {

		contender := domain.Contender{
			ContestID: contestID,
			Ownership: domain.OwnershipData{
				OrganizerID: contest.Ownership.OrganizerID,
			},
			RegistrationCode: uc.RegistrationCodeGenerator.Generate(registrationCodeLength),
		}

		if contender, err = uc.Repo.StoreContender(ctx, tx, contender); err != nil {
			tx.Rollback()
			return nil, errors.Wrap(err, 0)
		}

		contenders = append(contenders, contender)
	}

	tx.Commit()

	return contenders, err
}
