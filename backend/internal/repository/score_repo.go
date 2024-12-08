package repository

import (
	"context"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
	"gorm.io/gorm"
)

type scoreRecord struct {
	ContenderID int `gorm:"primaryKey"`
	Timestamp   time.Time
	Score       int
	Placement   int
	Finalist    bool
	Rankorder   int `gorm:"column:rank_order"`
}

func (scoreRecord) TableName() string {
	return "score"
}

func (r scoreRecord) fromDomain(score domain.Score) scoreRecord {
	return scoreRecord{
		ContenderID: int(score.ContenderID),
		Timestamp:   score.Timestamp,
		Score:       score.Score,
		Placement:   score.Placement,
		Finalist:    score.Finalist,
		Rankorder:   score.RankOrder,
	}
}

func (r *scoreRecord) toDomain() domain.Score {
	return domain.Score{
		Timestamp:   r.Timestamp,
		ContenderID: domain.ContenderID(r.ContenderID),
		Score:       r.Score,
		Placement:   r.Placement,
		Finalist:    r.Finalist,
		RankOrder:   r.Rankorder,
	}
}

func (d *Database) StoreScore(ctx context.Context, tx domain.Transaction, score domain.Score) (domain.Score, error) {
	var err error
	var record scoreRecord = scoreRecord{}.fromDomain(score)

	err = d.tx(tx).WithContext(ctx).Save(&record).Error
	switch err {
	case nil:
	case gorm.ErrForeignKeyViolated:
		return domain.Score{}, errors.New(domain.ErrNotFound)
	default:
		return domain.Score{}, errors.Wrap(err, 0)
	}

	return record.toDomain(), nil
}
