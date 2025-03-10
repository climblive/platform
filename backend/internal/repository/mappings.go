package repository

import (
	"database/sql"
	"time"

	"github.com/climblive/platform/backend/internal/database"
	"github.com/climblive/platform/backend/internal/domain"
)

func contenderToDomain(record database.GetContenderRow) domain.Contender {
	contender := domain.Contender{
		ID: domain.ContenderID(record.Contender.ID),
		Ownership: domain.OwnershipData{
			OrganizerID: domain.OrganizerID(record.Contender.OrganizerID),
			ContenderID: nillableIntToResourceID[domain.ContenderID](&record.Contender.ID),
		},
		ContestID:           domain.ContestID(record.Contender.ContestID),
		CompClassID:         domain.CompClassID(record.Contender.ClassID.Int32),
		RegistrationCode:    record.Contender.RegistrationCode,
		Name:                record.Contender.Name.String,
		PublicName:          record.Contender.Name.String,
		ClubName:            record.Contender.Club.String,
		Entered:             record.Contender.Entered.Time,
		WithdrawnFromFinals: record.Contender.WithdrawnFromFinals,
		Disqualified:        record.Contender.Disqualified,
	}

	if record.ContenderID.Valid {
		score := domain.Score{
			Timestamp:   record.Timestamp.Time,
			ContenderID: domain.ContenderID(record.ContenderID.Int32),
			Score:       int(record.Score.Int32),
			Placement:   int(record.Placement.Int32),
			Finalist:    record.Finalist.Bool,
			RankOrder:   int(record.RankOrder.Int32),
		}

		contender.Score = &score
	}

	return contender
}

func compClassToDomain(record database.CompClass) domain.CompClass {
	return domain.CompClass{
		ID: domain.CompClassID(record.ID),
		Ownership: domain.OwnershipData{
			OrganizerID: domain.OrganizerID(record.OrganizerID),
		},
		ContestID:   domain.ContestID(record.ContestID),
		Name:        record.Name,
		Description: record.Description.String,
		Color:       domain.ColorRGB(record.Color.String),
		TimeBegin:   record.TimeBegin,
		TimeEnd:     record.TimeEnd,
	}
}

func contestToDomain(record database.Contest) domain.Contest {
	contest := domain.Contest{
		ID: domain.ContestID(record.ID),
		Ownership: domain.OwnershipData{
			OrganizerID: domain.OrganizerID(record.OrganizerID),
		},
		Location:           record.Location.String,
		SeriesID:           domain.SeriesID(record.SeriesID.Int32),
		Protected:          record.Protected,
		Name:               record.Name,
		Description:        record.Description.String,
		FinalsEnabled:      record.FinalEnabled,
		QualifyingProblems: int(record.QualifyingProblems),
		Finalists:          int(record.Finalists),
		Rules:              record.Rules.String,
		GracePeriod:        time.Duration(record.GracePeriod) * time.Minute,
	}

	if !contest.FinalsEnabled {
		contest.Finalists = 0
	}

	return contest
}

func problemToDomain(record database.Problem) domain.Problem {
	return domain.Problem{
		ID: domain.ProblemID(record.ID),
		Ownership: domain.OwnershipData{
			OrganizerID: domain.OrganizerID(record.OrganizerID),
		},
		ContestID:          domain.ContestID(record.ContestID),
		Number:             int(record.Number),
		HoldColorPrimary:   record.HoldColorPrimary,
		HoldColorSecondary: record.HoldColorSecondary.String,
		Name:               record.Name.String,
		Description:        record.Description.String,
		PointsTop:          int(record.Points),
		PointsZone:         0,
		FlashBonus:         int(record.FlashBonus.Int32),
	}
}

func tickToDomain(record database.Tick) domain.Tick {
	attempts := func(isFlash bool) int {
		if isFlash {
			return 1
		}

		return 999
	}

	return domain.Tick{
		ID: domain.TickID(record.ID),
		Ownership: domain.OwnershipData{
			OrganizerID: domain.OrganizerID(record.OrganizerID),
			ContenderID: nillableIntToResourceID[domain.ContenderID](&record.ContenderID),
		},
		Timestamp:    record.Timestamp,
		ContestID:    domain.ContestID(record.ContestID),
		ProblemID:    domain.ProblemID(record.ProblemID),
		Top:          true,
		AttemptsTop:  attempts(record.Flash),
		Zone:         true,
		AttemptsZone: attempts(record.Flash),
	}
}

func userToDomain(record database.User) domain.User {
	return domain.User{
		ID:         domain.UserID(record.ID),
		Name:       record.Name,
		Username:   record.Username,
		Admin:      record.Admin,
		Organizers: make([]domain.OrganizerID, 0),
	}
}

func organizerToDomain(record database.Organizer) domain.Organizer {
	return domain.Organizer{
		ID: domain.OrganizerID(record.ID),
		Ownership: domain.OwnershipData{
			OrganizerID: domain.OrganizerID(record.ID),
		},
		Name:     record.Name,
		Homepage: record.Homepage.String,
	}
}

func makeNullString(value string) sql.NullString {
	if value == "" {
		return sql.NullString{}
	}

	return sql.NullString{
		Valid:  true,
		String: value,
	}
}

func makeNullInt32(value int32) sql.NullInt32 {
	if value == 0 {
		return sql.NullInt32{}
	}

	return sql.NullInt32{
		Valid: true,
		Int32: value,
	}
}

func makeNullTime(value time.Time) sql.NullTime {
	if value.IsZero() {
		return sql.NullTime{}
	}

	return sql.NullTime{
		Valid: true,
		Time:  value,
	}
}
