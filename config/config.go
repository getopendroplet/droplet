package config

import (
	"droplet/utils"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Config holds configs
type Config struct {
	DefaultRemote  string                                `yaml:"default-remote"`
	Remotes        map[string]Remote                     `yaml:"remotes"`
	Configs        map[string]string                     `yaml:"configs"`
	UserAgent      string                                `yaml:"-"`
	PromptPassword func(filename string) (string, error) `yaml:"-"`
}

// NewConfig returns a Config, optionally using default.
func NewConfig(defaults bool) *Config {
	config := &Config{}
	if defaults {
		config.DefaultRemote = "origin"
		// config.Remotes = map[string]Remote{
		// 	"origin": Remote{
		// 		Addr:     "https://",
		// 		Protocol: "http",
		// 	},
		// }
		config.Configs = map[string]string{
			"author":                          "Evaldas Leliuga",
			"author_email":                    "getopendroplet@gmail.com",
			"package_manager":                 "apk",
			"package_manager_action_by_stage": "true",
			"generator":                       "local",
		}
	}

	return config
}

// LoadConfig reads the configuration from the config file; if the file does
// not exist, it returns a default configuration.
func LoadConfig(name string) (*Config, error) {
	// Open the config file
	content, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("Unable to read the configuration file: %v", err)
	}

	// Decode the YAML document
	c := NewConfig(false)
	if err := yaml.Unmarshal(content, &c); err != nil {
		return nil, fmt.Errorf("Unable to decode the configuration: %v", err)
	}

	return c, nil
}

// SaveConfig writes the provided configuration to the config file.
func (c *Config) SaveConfig(name string) error {
	dir, _ := filepath.Split(name)
	if !utils.PathExists(dir) {
		if err := os.MkdirAll(dir, 0750); err != nil {
			return fmt.Errorf("Unable to create the configuration dir: %v", dir)
		}
	}

	// Create the config file (or truncate an existing one)
	f, err := os.Create(name)
	if err != nil {
		return fmt.Errorf("Unable to create the configuration file: %v", err)
	}
	defer f.Close()

	// Write the new config
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("Unable to marshal the configuration: %v", err)
	}

	if _, err := f.Write(data); err != nil {
		return fmt.Errorf("Unable to write the configuration: %v", err)
	}

	return nil
}
