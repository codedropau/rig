package config

import stringutils "github.com/codedropau/rig/internal/utils/string"

// Service a service declared in a Docker Compose file.
type Service struct {
	Labels      map[string]string `yaml:"labels"`
	Image       string            `yaml:"image"`
	Volumes     []string          `yaml:"volumes"`
	Environment []string          `yaml:"environment"`
	ExtraHosts  []string          `yaml:"extra_hosts"`
}

// MountsVolume will check if a volume is configured to be mounted for this Service.
func (s Service) MountsVolume(name string) bool {
	for _, volume := range s.Volumes {
		source, _, err := stringutils.SplitBySeparator(volume, ":")
		if err != nil {
			continue
		}

		if source == name {
			return true
		}
	}

	return false
}