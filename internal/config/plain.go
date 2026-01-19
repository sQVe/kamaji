package config

import (
	"os"
	"sync"
	"sync/atomic"
)

var (
	plainMode     bool
	plainOnce     sync.Once
	plainOverride atomic.Pointer[bool]
)

// IsPlain returns true when terminal styling should be disabled.
// Checks KAMAJI_PLAIN env var or NO_COLOR standard.
func IsPlain() bool {
	if v := plainOverride.Load(); v != nil {
		return *v
	}
	plainOnce.Do(func() {
		if isTruthy(os.Getenv("KAMAJI_PLAIN")) {
			plainMode = true
			return
		}
		if os.Getenv("NO_COLOR") != "" {
			plainMode = true
		}
	})
	return plainMode
}

func isTruthy(s string) bool {
	switch s {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}

// SetPlain forces plain mode (for testing).
func SetPlain(v bool) {
	plainOverride.Store(&v)
}

// ResetPlain clears the override and resets detection.
// Not thread-safe: only call from test setup, not concurrently with IsPlain.
func ResetPlain() {
	plainOverride.Store(nil)
	plainOnce = sync.Once{}
	plainMode = false
}
