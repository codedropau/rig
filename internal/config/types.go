package config

import (
	"time"
)

// Config used by Rig to provision environments.
type Config struct {
	Project     string             `yaml:"project"`
	Dockerfiles []string           `yaml:"dockerfiles"`
	Services    map[string]Service `yaml:"services"`
	Ingress     Ingress            `yaml:"ingress"`
	Volume      Volume             `yaml:"volume"`
	Retention   time.Duration      `yaml:"retention"`
}

// Service which will be used during a snapshot and run.
type Service struct {
	CPU    string `yaml:"cpu"`
	Memory string `yaml:"memory"`
}

// Ingress configuration.
type Ingress struct {
	Port int `yaml:"port"`
}

type Volume struct {
	From  string `yaml:"from"`
	User  string `yaml:"user"`
	Group string `yaml:"group"`
}
