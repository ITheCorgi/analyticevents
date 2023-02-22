package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ITheCorgi/analyticevents/internal/entity"
	"github.com/ITheCorgi/analyticevents/pkg/server"
	"github.com/gin-gonic/gin"
)

type EventWriter interface {
	// ConsumeEvents fetch line-divided jsons into batch, then stores into clickhouse
	ConsumeEvents(ctx context.Context, stream chan entity.AnalyticEvent)
}

func (h handler) AddEvents(ctx *gin.Context, params server.AddEventsParams) {
	ch := make(chan entity.AnalyticEvent, 30)
	defer close(ch)

	go h.consumer.ConsumeEvents(ctx, ch)

	ip := ctx.ClientIP()

	dec := json.NewDecoder(ctx.Request.Body)
	for dec.More() {
		select {
		case <-ctx.Done():
			fmt.Println("context is cancelled")
		default:
			var e entity.AnalyticEvent

			err := dec.Decode(&e)
			if err != nil {
				fmt.Println(err)
			}

			e.Time, _ = time.Parse("2020-12-01 23:59:00", e.ClientTime)

			e.ClientIP = ip
			e.ServerTime = time.Now().Local()

			ch <- e
		}
	}
}
