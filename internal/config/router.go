package config

type RouterConfig struct {
	PORT         string
	REGISTRY_URL string
}

func NewRouter() *RouterConfig {
	return &RouterConfig{
		PORT:         getString("PORT", "8081"),
		REGISTRY_URL: getString("REGISTRY_URL", "localhost:8010"),
	}
}
