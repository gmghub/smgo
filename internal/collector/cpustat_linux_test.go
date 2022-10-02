//go:build linux
// +build linux

package collector

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGetCPUStat(t *testing.T) {
	t.Run("sysstat returns data", func(t *testing.T) {
		now := time.Now()
		s := GetCPUStat()
		require.True(t, s.Time.After(now), "cpustat time should be actual")
		require.Greater(t, s.UserMode+s.SysMode+s.Idle, float32(0), "cpu load sum should be greater then 0: ", s)
	})
}
