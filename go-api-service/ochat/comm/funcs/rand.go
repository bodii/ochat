package funcs

import (
	"math/rand"
	"time"
)

// gets a Rand handle for random nanoseconds
//
// param:
//
// return:
//   - [*rand.Rand] rand handle
func GetUnixNanoRandSeed() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}
