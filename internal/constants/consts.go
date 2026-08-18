package constants

import "time"

const (
	// Registry constants
	StaleInterval = 15 * time.Second
	SweeperTick   = 5 * time.Second

	// Shard constants
	ReporterTick = 5 * time.Second
	NumBuckets   = 256 // should be power of 2

	// Poller constants
	PollerTick = 5 * time.Second

	// Dispatcher call timeouts per-method
	PutTimeout    = 2 * time.Second
	GetTimeout    = 2 * time.Second
	DeleteTimeout = 2 * time.Second
	ExistsTimeout = 2 * time.Second
	InfoTimeout   = 5 * time.Second

	// Common constants
	ClientTimeout = 2 * time.Second
)
