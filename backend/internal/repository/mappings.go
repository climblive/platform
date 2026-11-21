package repository

import (
	"database/sql"
	"time"

	"github.com/climblive/platform/backend/internal/database"
	"github.com/climblive/platform/backend/internal/domain"
	"uuid"
)

func competitorToDomain(record database.GetCompetitorRow) domain.Competitor {
	competitor := domain.Competitor{
		ID: domain.CompetitorID(record.Competitor.ID),
		Ownership: domain.OwnershipData{
			OrganizerID:  domain.OrganizerID(record.Competitor.OrganizerID),
			CompetitorID: nillableIntToResourceID[domain.CompetitorID](&record.Competitor.ID),
		},
		Result:              nil,
		CompetitionID:       domain.CompetitionID(record.Competitor.CompetitionID),
		CategoryID:          domain.CategoryID(record.Competitor.CategoryID.Int32),
		RegistrationCode:    record.Competitor.RegistrationCode,
		Name:                record.Competitor.Name.String,
		Entered:             record.Competitor.Entered.Time,
		WithdrawnFromFinals: record.Competitor.WithdrawnFromFinals,
		Disqualified:        record.Competitor.Disqualified,
		ScrubbedAt:          record.Competitor.ScrubbedAt.Time,
		ScrubBefore:         record.Competitor.ScrubBefore.Time,
	}

	if record.CompetitorID.Valid {
		result := domain.Result{
			Timestamp:    record.Timestamp.Time,
			CompetitorID: domain.CompetitorID(record.CompetitorID.Int32),
			Score:        record.Score.String,
			Placement:    int(record.Placement.Int32),
			Finalist:     record.Finalist.Bool,
			RankOrder:    int(record.RankOrder.Int32),
		}

		competitor.Result = &result
	}

	return competitor
}

func categoryToDomain(record database.Category) domain.Category {
	return domain.Category{
		ID: domain.CategoryID(record.ID),
		Ownership: domain.OwnershipData{
			OrganizerID:  domain.OrganizerID(record.OrganizerID),
			CompetitorID: nil,
		},
		CompetitionID: domain.CompetitionID(record.CompetitionID),
		Name:          record.Name,
		Description:   record.Description.String,
		MinAge:        int(record.MinAge),
		MaxAge:        int(record.MaxAge),
		TimeBegin:     record.TimeBegin,
		TimeEnd:       record.TimeEnd,
	}
}

func competitionToDomain(record database.Competition) domain.Competition {
	competition := domain.Competition{
		ID: domain.CompetitionID(record.ID),
		Ownership: domain.OwnershipData{
			OrganizerID:  domain.OrganizerID(record.OrganizerID),
			CompetitorID: nil,
		},
		TimeBegin:             time.Time{},
		TimeEnd:               time.Time{},
		RegisteredCompetitors: 0,
		ArchivedAt:            record.ArchivedAt.Time,
		Location:              record.Location.String,
		Country:               record.Country,
		SeriesID:              domain.SeriesID(record.SeriesID.Int32),
		Name:                  record.Name,
		Description:           record.Description.String,
		QualifyingProblems:    int(record.QualifyingProblems),
		Finalists:             int(record.Finalists),
		UsePoints:             record.UsePoints,
		PooledPoints:          record.PooledPoints,
		MaxAttempts:           int(record.MaxAttempts),
		PointDeduction:        domain.TenthsFromFloat64(record.PointDeduction),
		Info:                  record.Info.String,
		GracePeriod:           time.Duration(record.GracePeriod) * time.Minute,
		NameRetentionTime:     time.Duration(record.NameRetentionTime) * time.Minute,
		Created:               record.Created,
	}

	return competition
}

func problemToDomain(record database.Problem) domain.Problem {
	return domain.Problem{
		ID: domain.ProblemID(record.ID),
		Ownership: domain.OwnershipData{
			OrganizerID:  domain.OrganizerID(record.OrganizerID),
			CompetitorID: nil,
		},
		CompetitionID:      domain.CompetitionID(record.CompetitionID),
		Number:             int(record.Number),
		HoldColorPrimary:   record.HoldColorPrimary,
		HoldColorSecondary: record.HoldColorSecondary.String,
		Description:        record.Description.String,
		Zone1Enabled:       record.Zone1Enabled,
		Zone2Enabled:       record.Zone2Enabled,
		ProblemValue: domain.ProblemValue{
			PointsZone1: int(record.PointsZone1.Int32),
			PointsZone2: int(record.PointsZone2.Int32),
			PointsTop:   int(record.PointsTop),
			FlashBonus:  int(record.FlashBonus.Int32),
		},
	}
}

func tickToDomain(record database.Tick) domain.Tick {
	return domain.Tick{
		ID: domain.TickID(record.ID),
		Ownership: domain.OwnershipData{
			OrganizerID:  domain.OrganizerID(record.OrganizerID),
			CompetitorID: nillableIntToResourceID[domain.CompetitorID](&record.CompetitorID),
		},
		Timestamp:     record.Timestamp,
		CompetitionID: domain.CompetitionID(record.CompetitionID),
		Revision:      int(record.Revision),
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
			OrganizerID:  domain.OrganizerID(record.ID),
			CompetitorID: nil,
		},
		Name: record.Name,
	}
}

func raffleToDomain(record database.Raffle) domain.Raffle {
	return domain.Raffle{
		ID: domain.RaffleID(record.ID),
		Ownership: domain.OwnershipData{
			OrganizerID:  domain.OrganizerID(record.OrganizerID),
			CompetitorID: nil,
		},
		CompetitionID: domain.CompetitionID(record.CompetitionID),
	}
}

func raffleWinnerToDomain(record database.RaffleWinner, name string, scrubbedAt time.Time) domain.RaffleWinner {
	return domain.RaffleWinner{
		ID: domain.RaffleWinnerID(record.ID),
		Ownership: domain.OwnershipData{
			OrganizerID:  domain.OrganizerID(record.OrganizerID),
			CompetitorID: nil,
		},
		RaffleID:             domain.RaffleID(record.RaffleID),
		CompetitorID:         domain.CompetitorID(record.CompetitorID),
		CompetitorName:       name,
		CompetitorScrubbedAt: scrubbedAt,
		Timestamp:            record.Timestamp,
	}
}

func organizerInviteToDomain(record database.OrganizerInvite, organizerName string) domain.OrganizerInvite {
	return domain.OrganizerInvite{
		ID:            domain.OrganizerInviteID(uuid.MustParse(record.ID)),
		OrganizerID:   domain.OrganizerID(record.OrganizerID),
		OrganizerName: organizerName,
		ExpiresAt:     record.ExpiresAt,
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
