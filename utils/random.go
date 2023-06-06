package utils

import (
	"math/rand"
	"time"
)

func GetRandomWithAll(min, max int) int64 {
	rand.Seed(time.Now().UnixNano())
	return int64(rand.Intn(max-min+1) + min)
}
