package config

// Config is an object which encapsulates a Docker Compose file.
type Config struct {
	Services map[string]*Service `yaml:"services"`
	Volumes  map[string]*Volume  `yaml:"volumes"`
}

// Volume which is mounted into Services.
type Volume struct {
	Driver string     `yaml:"driver"`
	Opts   VolumeOpts `yaml:"driver_opts"`
}

// VolumeOpts which are used when mounting Volumes.
type VolumeOpts struct {
	Device string `yaml:"device"`
}
