package domain

type Problem struct {
	ID                 ResourceID
	ContestID          ResourceID
	Number             int
	HoldColorPrimary   string
	HoldColorSecondary string
	Name               string
	Description        string
	Points             int
	FlashBonus         int
}

type ProblemUsecase interface {
}
