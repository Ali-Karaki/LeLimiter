package utils

import (
	"math/rand"
	"rate-limiter/pkg/constants"
)

func GetRandomReqCount() int {
	diff := constants.ConcReqsCountMinMax.Max - constants.ConcReqsCountMinMax.Min + 1
	return rand.Intn(diff) + constants.ConcReqsCountMinMax.Min
}

func GetRandomUserID(users []string) string {
	return users[rand.Intn(len(users))]
}
