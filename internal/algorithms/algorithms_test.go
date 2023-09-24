package algorithms

import (
	"testing"
	"time"

	"github.com/leonid-tankov/OTUS_go_final_project/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestMultiArmedBandit(t *testing.T) {
	counters := []storage.Counter{
		{
			BannerID: 1,
			Count:    1,
		},
		{
			BannerID: 2,
			Count:    5,
		},
		{
			BannerID: 3,
			Count:    6,
		},
		{
			BannerID: 4,
			Count:    6,
		},
		{
			BannerID: 5,
			Count:    3,
		},
	}

	t.Run("probability - 1", func(t *testing.T) { // only max
		var ids []int
		for i := 0; i < 50; i++ {
			time.Sleep(1 * time.Nanosecond)
			id := MultiArmedBandit(counters, 1)
			if sliceContains(ids, int(id)) {
				continue
			}
			ids = append(ids, int(id))
		}
		require.Equal(t, 2, len(ids))
		require.Contains(t, ids, 3)
		require.Contains(t, ids, 4)
	})

	t.Run("probability - 0", func(t *testing.T) { // random
		var ids []int
		for i := 0; i < 50; i++ {
			time.Sleep(1 * time.Nanosecond)
			id := MultiArmedBandit(counters, 0)
			if sliceContains(ids, int(id)) {
				continue
			}
			ids = append(ids, int(id))
		}
		require.Equal(t, 5, len(ids))
	})
}

func sliceContains(slice []int, element int) bool {
	for _, a := range slice {
		if a == element {
			return true
		}
	}
	return false
}
