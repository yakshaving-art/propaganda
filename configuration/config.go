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
		return fmt.Errorf("failed to parse: %s", err)
	}

	config = c

	return nil
}

// GetConfiguration returns the current configuration
func GetConfiguration() Configuration {
	return config
}

// GetChannel returns the specific channel for the repo, or the default one if there is no specific set.
func (c Configuration) GetChannel(repoFullName string) string {
	channel, ok := c.Repositories[repoFullName]
	if ok {
		return channel
	}
	return c.DefaultChannel
}
