package config

type ShardConfig struct {
	PORT         string
	REGISTRY_URL string
	SHARD_ADDR   string
}

func NewShard() *ShardConfig {
	cfg := &ShardConfig{
		PORT:         getString("PORT", "8080"),
		REGISTRY_URL: getString("REGISTRY_URL", "localhost:8010"),
		SHARD_ADDR:   getString("SHARD_ADDR", ""),
	}

	if cfg.SHARD_ADDR == "" {
		cfg.SHARD_ADDR = "localhost:" + cfg.PORT
	}

	cfg.PORT = ":" + cfg.PORT

	return cfg
}
