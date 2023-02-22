package repository

import (
	"context"
	"errors"

	"github.com/ITheCorgi/analyticevents/internal/entity"
	"github.com/uptrace/go-clickhouse/ch"
)

type events struct {
	db *ch.DB
}

func New(db *ch.DB) events {
	return events{
		db: db,
	}
}

func (e events) SaveBatch(ctx context.Context, batch []entity.AnalyticEvent) error {
	res, err := e.db.NewInsert().Model(&batch).Exec(ctx)
	if err != nil {
		return err
	}

	act, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if act != int64(len(batch)) {
		return errors.New("mismatched inserted and expected event records into clickhouse db")
	}

	return nil
}
