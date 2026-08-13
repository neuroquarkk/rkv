package constants

import "time"

const (
	// Registry constants
	StaleInterval = 15 * time.Second
	SweeperTick   = 5 * time.Second

	// Shard constants
	ReporterTick = 5 * time.Second

	// Poller constants
	PollerTick = 5 * time.Second

	// Common constants
	ClientTimeout = 2 * time.Second
)
