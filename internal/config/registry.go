package config

type RegistryConfig struct {
	PORT string
}

func NewRegistry() *RegistryConfig {
	return &RegistryConfig{
		PORT: getString("PORT", "8010"),
	}
}
