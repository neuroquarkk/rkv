package config

import "os"

type ShardConfig struct {
	PORT         string
	REGISTRY_URL string
	SHARD_ADDR   string
}

func inContainer() bool {
	_, err := os.Stat("/.dockerenv")
	return err == nil
}

func resolveShardAddr(port string) string {
	if addr := getString("SHARD_ADDR", ""); addr != "" {
		return addr
	}
	if inContainer() {
		return os.Getenv("HOSTNAME") + ":" + port
	}
	return "localhost:" + port
}

func NewShard() *ShardConfig {
	cfg := &ShardConfig{
		PORT:         getString("PORT", "8080"),
		REGISTRY_URL: getString("REGISTRY_URL", "localhost:8010"),
	}

	cfg.SHARD_ADDR = resolveShardAddr(cfg.PORT)
	cfg.PORT = ":" + cfg.PORT

	return cfg
}
