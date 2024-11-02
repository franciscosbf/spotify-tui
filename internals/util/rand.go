package util

import (
	"math/rand"
	"time"
)

func NewRand() *rand.Rand {
	seed := time.Now().UnixNano()
	source := rand.NewSource(seed)

	return rand.New(source)
}
