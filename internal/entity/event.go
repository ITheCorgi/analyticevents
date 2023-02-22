package entity

import (
	"github.com/google/uuid"
	"github.com/uptrace/go-clickhouse/ch"
	"time"
)

type AnalyticEvent struct {
	ch.CHModel `ch:"analytics"`

	DeviceID   uuid.UUID `json:"device_id" ch:"device_id"`
	Time       time.Time `json:"-" ch:"client_time"`
	DeviceOS   string    `json:"device_os" ch:"device_os"`
	Session    string    `json:"session" ch:"session"`
	Sequence   uint64    `json:"sequence" ch:"sequence"`
	Event      string    `json:"event" ch:"event"`
	ParamInt   uint64    `json:"param_int" ch:"param_int"`
	ParamStr   string    `json:"param_str" ch:"param_str"`
	ClientIP   string    `json:"-" ch:"client_ip"`
	ClientTime string    `json:"client_time" ch:"-"`
	ServerTime time.Time `json:"-" ch:"server_time"`
}
