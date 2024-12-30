package repository

import (
	"database/sql"
	"time"

	"github.com/climblive/platform/backend/internal/database"
	"github.com/climblive/platform/backend/internal/domain"
)

func nullTimeToTime(time sql.NullTime) *time.Time {
	if time.Valid {
		return &time.Time
	}

	return nil
}

func contenderToDomain(contender database.Contender) domain.Contender {
	return domain.Contender{
		ID: domain.ContenderID(contender.ID),
		Ownership: domain.OwnershipData{
			OrganizerID: domain.OrganizerID(contender.OrganizerID),
			ContenderID: nillableIntToResourceID[domain.ContenderID](&contender.ID),
		},
		ContestID:           domain.ContestID(contender.ContestID),
		CompClassID:         domain.CompClassID(contender.ClassID.Int32),
		RegistrationCode:    contender.RegistrationCode,
		Name:                contender.Name.String,
		PublicName:          contender.Name.String,
		ClubName:            contender.Club.String,
		Entered:             nullTimeToTime(contender.Entered),
		WithdrawnFromFinals: contender.WithdrawnFromFinals,
		Disqualified:        contender.Disqualified,
	}
}

func scoreToDomain(score database.Score) domain.Score {
	return domain.Score{
		Timestamp:   score.Timestamp,
		ContenderID: domain.ContenderID(score.ContenderID),
		Score:       int(score.Score),
		Placement:   int(score.Placement),
		Finalist:    score.Finalist,
		RankOrder:   int(score.RankOrder),
	}
}

func compClassToDomain(compClass database.CompClass) domain.CompClass {
	return domain.CompClass{
		ID: domain.CompClassID(compClass.ID),
		Ownership: domain.OwnershipData{
			OrganizerID: domain.OrganizerID(compClass.OrganizerID),
		},
		ContestID:   domain.ContestID(compClass.ContestID),
		Name:        compClass.Name,
		Description: compClass.Description.String,
		Color:       domain.ColorRGB(compClass.Color.String),
		TimeBegin:   compClass.TimeBegin,
		TimeEnd:     compClass.TimeEnd,
	}
}

func contestToDomain(contest database.Contest) domain.Contest {
	return domain.Contest{
		ID: domain.ContestID(contest.ID),
		Ownership: domain.OwnershipData{
			OrganizerID: domain.OrganizerID(contest.OrganizerID),
		},
		Location:           contest.Location.String,
		SeriesID:           domain.SeriesID(contest.SeriesID.Int32),
		Protected:          contest.Protected,
		Name:               contest.Name,
		Description:        contest.Description.String,
		FinalsEnabled:      contest.FinalEnabled,
		QualifyingProblems: int(contest.QualifyingProblems),
		Finalists:          int(contest.Finalists),
		Rules:              contest.Rules.String,
		GracePeriod:        time.Duration(contest.GracePeriod),
	}
}

func problemToDomain(problem database.Problem) domain.Problem {
	return domain.Problem{
		ID: domain.ProblemID(problem.ID),
		Ownership: domain.OwnershipData{
			OrganizerID: domain.OrganizerID(problem.OrganizerID),
		},
		ContestID:          domain.ContestID(problem.ContestID),
		Number:             int(problem.Number),
		HoldColorPrimary:   problem.HoldColorPrimary,
		HoldColorSecondary: problem.HoldColorSecondary.String,
		Name:               problem.Name.String,
		Description:        problem.Description.String,
		PointsTop:          int(problem.Points),
		PointsZone:         0,
		FlashBonus:         int(problem.FlashBonus.Int32),
	}
}

func tickToDomain(tick database.Tick) domain.Tick {
	attempts := func(isFlash bool) int {
		if isFlash {
			return 1
		}

		return 999
	}

	return domain.Tick{
		ID: domain.TickID(tick.ID),
		Ownership: domain.OwnershipData{
			OrganizerID: domain.OrganizerID(tick.OrganizerID),
			ContenderID: nillableIntToResourceID[domain.ContenderID](&tick.ContenderID),
		},
		Timestamp:    tick.Timestamp,
		ContestID:    domain.ContestID(tick.ContestID),
		ProblemID:    domain.ProblemID(tick.ProblemID),
		Top:          true,
		AttemptsTop:  attempts(tick.Flash),
		Zone:         true,
		AttemptsZone: attempts(tick.Flash),
	}
}
