package algorithms

import (
	"math/rand"
	"time"

	"github.com/leonid-tankov/OTUS_go_final_project/internal/storage"
)

func MultiArmedBandit(counters []storage.Counter, probability float64) int64 {
	rand.Seed(time.Now().UnixNano())
	epsilon := rand.Float64() //nolint:gosec
	if probability > epsilon {
		return randMax(counters)
	}
	return counters[rand.Intn(len(counters))].BannerID //nolint:gosec
}

func randMax(counters []storage.Counter) int64 {
	var max int
	for _, counter := range counters {
		if counter.Count > max {
			max = counter.Count
		}
	}
	var maxIDs []int64
	for _, counter := range counters {
		if counter.Count == max {
			maxIDs = append(maxIDs, counter.BannerID)
		}
	}
	rand.Seed(time.Now().UnixNano())
	return maxIDs[rand.Intn(len(maxIDs))] //nolint:gosec
}
