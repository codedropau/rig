package config

import (
	"errors"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Load the Docker Compose files.
func Load(filePaths []string) (*Config, error) {
	if len(filePaths) == 0 {
		return nil, errors.New("config file not provided")
	}

	if len(filePaths) == 1 {
		return load(filePaths[0])
	}

	return loadMultiplle(filePaths[0], filePaths[1:])
}

func load(primary string) (*Config, error) {
	file, err := ioutil.ReadFile(primary)
	if err != nil {
		return nil, err
	}

	config := &Config{}

	err = yaml.Unmarshal(file, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func loadMultiplle(primary string, overrides []string) (*Config, error) {
	config, err := load(primary)
	if err != nil {
		return nil, err
	}

	for _, override := range overrides {
		o, err := load(override)
		if err != nil {
			return nil, err
		}

		err = merge(config, *o)
		if err != nil {
			return nil, err
		}
	}

	return config, nil
}