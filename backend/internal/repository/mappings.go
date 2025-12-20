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
		Archived:           record.Archived,
		Location:           record.Location.String,
		SeriesID:           domain.SeriesID(record.SeriesID.Int32),
		Name:               record.Name,
		Description:        record.Description.String,
		QualifyingProblems: int(record.QualifyingProblems),
		Finalists:          int(record.Finalists),
		Rules:              record.Rules.String,
		GracePeriod:        time.Duration(record.GracePeriod) * time.Minute,
	}

	if !record.FinalEnabled {
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
		Description:        record.Description.String,
		Zone1Enabled:       record.Zone1Enabled,
		Zone2Enabled:       record.Zone2Enabled,
		PointsZone1:        int(record.PointsZone1.Int32),
		PointsZone2:        int(record.PointsZone2.Int32),
		PointsTop:          int(record.PointsTop),
		FlashBonus:         int(record.FlashBonus.Int32),
	}
}

func tickToDomain(record database.Tick) domain.Tick {
	return domain.Tick{
		ID: domain.TickID(record.ID),
		Ownership: domain.OwnershipData{
			OrganizerID: domain.OrganizerID(record.OrganizerID),
			ContenderID: nillableIntToResourceID[domain.ContenderID](&record.ContenderID),
		},
		Timestamp:     record.Timestamp,
		ContestID:     domain.ContestID(record.ContestID),
		ProblemID:     domain.ProblemID(record.ProblemID),
		Zone1:         record.Zone1,
		AttemptsZone1: int(record.AttemptsZone1),
		Zone2:         record.Zone2,
		AttemptsZone2: int(record.AttemptsZone2),
		Top:           record.Top,
		AttemptsTop:   int(record.AttemptsTop),
	}
}

func userToDomain(record database.User) domain.User {
	return domain.User{
		ID:         domain.UserID(record.ID),
		Username:   record.Username,
		Admin:      record.Admin,
		Organizers: make([]domain.Organizer, 0),
	}
}

func organizerToDomain(record database.Organizer) domain.Organizer {
	return domain.Organizer{
		ID: domain.OrganizerID(record.ID),
		Ownership: domain.OwnershipData{
			OrganizerID: domain.OrganizerID(record.ID),
		},
		Name: record.Name,
	}
}

func raffleToDomain(record database.Raffle) domain.Raffle {
	return domain.Raffle{
		ID: domain.RaffleID(record.ID),
		Ownership: domain.OwnershipData{
			OrganizerID: domain.OrganizerID(record.OrganizerID),
		},
		ContestID: domain.ContestID(record.ContestID),
	}
}

func raffleWinnerToDomain(record database.RaffleWinner, name string) domain.RaffleWinner {
	return domain.RaffleWinner{
		ID:            domain.RaffleWinnerID(record.ID),
		Ownership:     domain.OwnershipData{OrganizerID: domain.OrganizerID(record.OrganizerID)},
		RaffleID:      domain.RaffleID(record.RaffleID),
		ContenderID:   domain.ContenderID(record.ContenderID),
		ContenderName: name,
		Timestamp:     record.Timestamp,
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
