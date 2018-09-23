package configuration

import (
	"fmt"
	"gopkg.in/yaml.v2"
)

// Configuration holds the general configuration for the executable
type Configuration struct {
	DefaultChannel string            `yaml:"default_channel"`
	Repositories   map[string]string `yaml:"repositories"`
}

var config Configuration

// Load loads the configuration and sets it globally if loading succeeded
func Load(in []byte) error {
	c := Configuration{}
	err := yaml.Unmarshal(in, &c)
	if err != nil {
		return fmt.Errorf("failed to parse configuration: %s", err)
	}

	config = c

	return nil
}

// GetConfiguration returns the current configuration
func GetConfiguration() Configuration {
	return config
}
