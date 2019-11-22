package config

// Config is an object which encapsulates a Docker Compose file.
type Config struct {
	Services map[string]*Service `yaml:"services"`
	Volumes  map[string]*Volume `yaml:"volumes"`
}

// Service a service declared in a Docker Compose file.
type Service struct {
	Labels      map[string]string `yaml:"labels"`
	Image       string            `yaml:"image"`
	Volumes     []string           `yaml:"volumes"`
	Ports       []string            `yaml:"ports"`
	Environment []string      `yaml:"environment"`
	ExtraHosts  []string        `yaml:"extra_hosts"`
}

type Volume struct {
	Driver string `yaml:"driver"`
	Opts VolumeOpts `yaml:"driver_opts"`
}

type VolumeOpts struct {
	Device string `yaml:"device"`
}
