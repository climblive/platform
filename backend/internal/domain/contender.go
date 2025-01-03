package domain

import (
	"context"
)

type ContenderUseCase interface {
	GetContender(ctx context.Context, contenderID ContenderID) (Contender, error)
	GetContenderByCode(ctx context.Context, registrationCode string) (Contender, error)
	GetContendersByCompClass(ctx context.Context, compClassID CompClassID) ([]Contender, error)
	GetContendersByContest(ctx context.Context, contestID ContestID) ([]Contender, error)
	PatchContender(ctx context.Context, contenderID ContenderID, patch ContenderPatch) (Contender, error)
	DeleteContender(ctx context.Context, contenderID ContenderID) error
	CreateContenders(ctx context.Context, contestID ContestID, number int) ([]Contender, error)
}

type CodeGenerator interface {
	Generate(length int) string
}
