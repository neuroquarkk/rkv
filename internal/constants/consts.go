package constants

import "time"

// Regsitry Constants
const (
	// member considered dead if no heartbeat within this window
	StaleInterval = 15 * time.Second

	// how often regsitry scans state for stale members
	SweeperTick = 5 * time.Second
)

// Shard Constants
const (
	// how often a shard sends a heartbeat to the registry
	ReporterTick = 5 * time.Second

	// lock striping fanout. should be power of 2
	NumBuckets = 256

	// keys sampled per eviction pass
	EvictionSampleSize = 5

	// total store budget across all budget (only values are taken into consideration)
	MaxStorageSize = 512 * 1024 * 1024 // 512mb => 2mb/bucket
)

// Poller Constants
const (
	// how often the router polls the registry for members list
	PollerTick = 5 * time.Second
)

// Dispatcher call timeouts per-method
const (
	PutTimeout    = 2 * time.Second
	GetTimeout    = 2 * time.Second
	DeleteTimeout = 2 * time.Second
	ExistsTimeout = 2 * time.Second
	InfoTimeout   = 5 * time.Second
)

// Common constants
const (
	ClientTimeout = 2 * time.Second

	// Key size bytes
	MaxKeySize = 256
	// Value size bytes
	MaxValueSize = 1024
)
