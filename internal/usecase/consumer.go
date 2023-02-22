package usecase

import (
	"context"

	"github.com/ITheCorgi/analyticevents/internal/entity"
	"go.uber.org/zap"
)

type Writer interface {
	// SaveBatch save batch events into database
	SaveBatch(ctx context.Context, batch []entity.AnalyticEvent) error
}

type consumer struct {
	repo Writer
	log  *zap.Logger
}

func New(log *zap.Logger, eventRepo Writer) consumer {
	return consumer{
		repo: eventRepo,
		log:  log,
	}
}

func (c consumer) ConsumeEvents(ctx context.Context, stream chan entity.AnalyticEvent) {
	var batch []entity.AnalyticEvent

	for {
		select {
		case <-ctx.Done():
			c.log.Info("context is canceled, graceful exiting")

			if len(batch) > 0 {
				err := c.repo.SaveBatch(ctx, batch)
				if err != nil {
					c.log.Error("failed to save event batch into ch", zap.Error(err))
					return
				}
			}

			return

		case e, ok := <-stream:
			if !ok {
				c.log.Info("stream channel is closed, graceful exiting")

				err := c.repo.SaveBatch(ctx, batch)
				if err != nil {
					c.log.Error("failed to save event batch into ch", zap.Error(err))
					return
				}

				return
			}

			batch = append(batch, e)
		}
	}
}
