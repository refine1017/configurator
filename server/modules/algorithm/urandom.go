package algorithm

import "math/rand"

// 从[min,max]中产生一个数
func Random(min, max int64) int64 {
	if min > max {
		return 0
	}

	return rand.Int63n(max-min+1) + min
}