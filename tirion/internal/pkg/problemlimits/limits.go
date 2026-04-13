package problemlimits

import "fmt"

// Bounds for problem limits (aligned with the judge / DB schema).
const (
	MinTimeLimitMs   int64 = 100
	MaxTimeLimitMs   int64 = 600_000 // 10 minutes
	MinMemoryLimitKb int64 = 1024    // 1 MiB
	MaxMemoryLimitKb int64 = 524_288 // 512 MiB
)

// Validate returns an error if limits are outside allowed ranges.
func Validate(timeLimitMs, memoryLimitKb int64) error {
	if timeLimitMs < MinTimeLimitMs || timeLimitMs > MaxTimeLimitMs {
		return fmt.Errorf("time_limit_ms must be between %d and %d inclusive", MinTimeLimitMs, MaxTimeLimitMs)
	}
	if memoryLimitKb < MinMemoryLimitKb || memoryLimitKb > MaxMemoryLimitKb {
		return fmt.Errorf("memory_limit_kb must be between %d and %d inclusive", MinMemoryLimitKb, MaxMemoryLimitKb)
	}
	return nil
}
