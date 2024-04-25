package constants

import (
	"rate-limiter/api/types"
	"time"
)

const CtxKey types.CtxKey = "key"

const BaseURL string = "http://localhost:8080"

const MaxReqAllowedPerWindow int64 = 15

const WindowDuration int64 = 10 * 1000 // in milliseconds (1 minute)

var ConcReqsCountMinMax types.MinMax = types.MinMax{Min: 20, Max: 40} // go doesnt support const for complex types

const CallsTimeout time.Duration = 15 * 60 * time.Second // (5 minutes)