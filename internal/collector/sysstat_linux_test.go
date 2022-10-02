//go:build linux
// +build linux

package collector

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGetSysStat(t *testing.T) {
	t.Run("sysstat returns data", func(t *testing.T) {
		now := time.Now()
		s := GetSysStat()
		require.True(t, s.Time.After(now), "sysstat time should be actual")
		require.Greater(t, s.LoadAvg1m, float32(0), "sysstat loadavg1m should be greater then 0")
	})
}
