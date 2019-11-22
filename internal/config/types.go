package config

type Config struct {
	Dockerfiles []string `yaml:"dockerfiles"`
	Services []string `yaml:"services"`
}
