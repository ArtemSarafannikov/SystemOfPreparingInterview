package config

func (c *Config) normalizeJudgeDefaults() {
	if c.Judge.SandboxMemoryOverheadMB <= 0 {
		c.Judge.SandboxMemoryOverheadMB = 256
	}
	if c.Judge.MaxRiverWorkers <= 0 {
		c.Judge.MaxRiverWorkers = 10
	}
	if c.Judge.MaxConcurrentSandboxes <= 0 {
		c.Judge.MaxConcurrentSandboxes = c.Judge.MaxRiverWorkers
	}
}
