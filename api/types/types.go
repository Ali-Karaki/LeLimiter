package types

import (
	"github.com/redis/go-redis/v9"
)

type CtxKey string

type MinMax struct {
	Min int
	Max int
}

type Request struct {
	UserID string `json:"userId"`
}

type Response struct {
	Success           bool
	Reason 			  string
	Count             int64
	TimeBeforeNextReq string
}

type Controller struct {
	Redis *redis.Client
}
