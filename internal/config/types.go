package config

import "time"

// Config used by Rig to provision environments.
type Config struct {
	Project     string        `yaml:"project"`
	Dockerfiles []string      `yaml:"dockerfiles"`
	Services    []string      `yaml:"services"`
	Ingress     Ingress       `yaml:"ingress"`
	Volume Volume `yaml:"volume"`
	Retention   time.Duration `yaml:"retention"`
}

// Ingress configuration.
type Ingress struct {
	Port int `yaml:"port"`
}

type Volume struct {
	From string `yaml:"from"`
	User string `yaml:"user"`
	Group string `yaml:"group"`
}