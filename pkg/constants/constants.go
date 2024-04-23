package constants

import (
	"rate-limiter/api/types"
	"time"
)

const CtxKey types.CtxKey = "key"

const BaseURL string = "http://localhost:8080"

const MaxReqAllowedPerWindow int64 = 10

const WindowDuration int64 = 10 * 1000 // in milliseconds (1 minute)

var ConcReqsCountMinMax types.MinMax = types.MinMax{Min: 10, Max: 20} // go doesnt support const for complex types

const CallsTimeout time.Duration = 5 * 60 * time.Second // (10 minutes)