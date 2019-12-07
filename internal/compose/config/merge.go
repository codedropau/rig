package config

import stringutils "github.com/codedropau/rig/internal/utils/string"

// Helper function to merge configuration.
func merge(existing *Config, override Config) error {
	for name, service := range existing.Services {
		// Ensure the override has this service.
		if _, ok := override.Services[name]; !ok {
			continue
		}

		overrideService := override.Services[name]

		for labelName, labelValue := range overrideService.Labels {
			service.Labels[labelName] = labelValue
		}

		if overrideService.Image != "" {
			service.Image = overrideService.Image
		}

		environment, err := mergeWithKeyValue(service.Environment, overrideService.Environment, "=")
		if err != nil {
			return err
		}
		service.Environment = environment

		extra, err := mergeWithKeyValue(service.ExtraHosts, overrideService.ExtraHosts, ":")
		if err != nil {
			return err
		}
		service.ExtraHosts = extra

		service.Volumes = mergeIfMissing(service.Volumes, overrideService.Volumes)
	}

	for name, volume := range existing.Volumes {
		// Ensure the override has this service.
		if _, ok := override.Volumes[name]; !ok {
			continue
		}

		overrideVolume := override.Volumes[name]

		if overrideVolume.Driver != "" {
			volume.Driver = overrideVolume.Driver
		}

		if overrideVolume.Opts.Device != "" {
			volume.Opts.Device = overrideVolume.Opts.Device
		}
	}

	return nil
}

// Merge into a list and override existing if nessecary.
func mergeWithKeyValue(current, overrides []string, separator string) ([]string, error) {
	var merged []string

	for _, c := range current {
		name, _, err := stringutils.SplitBySeparator(c, separator)
		if err != nil {
			return merged, err
		}

		for _, override := range overrides {
			overrideName, _, err := stringutils.SplitBySeparator(override, separator)
			if err != nil {
				return merged, err
			}

			// We have an override we don't worry about including it on the list, it will be included later.
			if name == overrideName {
				continue
			}
		}

		merged = append(merged, c)
	}

	merged = append(merged, overrides...)

	return merged, nil
}

// Merge overrides into a current slice if they do not exist.
func mergeIfMissing(current, overrides []string) []string {
	for _, override := range overrides {
		if stringutils.Contains(current, override) {
			continue
		}

		current = append(current, override)
	}

	return current
}
