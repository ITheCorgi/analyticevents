package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/ITheCorgi/analyticevents/internal/config"
	"github.com/uptrace/go-clickhouse/ch"
)

func ProvideChConnection(ctx context.Context, cfg *config.Config) (*ch.DB, error) {
	conn := ch.Connect(
		ch.WithDSN(fmt.Sprintf("clickhouse://%s:%s/%s", cfg.DB.Host, cfg.DB.Port, cfg.DB.Database)),
		ch.WithUser(cfg.DB.Username),
		ch.WithPassword(cfg.DB.Password),
		ch.WithAutoCreateDatabase(true),
		ch.WithInsecure(true),
	)

	connCtx, cancel := context.WithTimeout(ctx, time.Second*20)
	defer cancel()

	err := conn.Ping(connCtx)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
