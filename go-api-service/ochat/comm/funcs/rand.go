package funcs

import (
	"math/rand"
	"time"
)

// 获取纳秒的随机种子
func GetUnixNanoRandSeed() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}
