package store

import (
	"math"
	"rkv/internal/constants"
)

func (b *bucket) evictLocked() {
	for b.size > b.budget {
		i := 0

		var oldestKey string
		var oldestTime int64 = math.MaxInt64
		for k, e := range b.data {
			if i >= constants.EvictionSampleSize {
				break
			}
			accessed := e.accessedAt.Load()
			if accessed < oldestTime {
				oldestTime = accessed
				oldestKey = k
			}
			i++
		}

		if oldestKey == "" {
			break
		}

		b.size -= uint64(len(b.data[oldestKey].value))
		delete(b.data, oldestKey)
	}
}
