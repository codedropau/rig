package config

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

const (
	// DefaultRoutingPort which is used when accessing the environment via Ingress.
	DefaultRoutingPort = 8080
	// DefaultDockerFile which is used when loading the projects configuration.
	DefaultDockerFile = "docker-compose.yml"
	// DefaultRetention will set how long an environment should be available before cleaned up.
	DefaultRetention = time.Hour * 120
	// DefaultVolumeFrom will declare the base image which stores the volume data.
	DefaultVolumeFrom = "alpine:latest"
	// DefaultVolumeUser permission which will be applied to the volume data.
	DefaultVolumeUser = "root"
	// DefaultVolumeGroup permission which will be applied to the volume data.
	DefaultVolumeGroup = "root"
)

// Load configuration file.
func Load(file string) (*Config, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	config := &Config{
		Dockerfiles: []string{
			DefaultDockerFile,
		},
		Ingress: Ingress{
			Port: DefaultRoutingPort,
		},
		Volume: Volume{
			From:  DefaultVolumeFrom,
			User:  DefaultVolumeUser,
			Group: DefaultVolumeGroup,
		},
		Retention: DefaultRetention,
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
