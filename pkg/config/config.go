package config

import (
	"bufio"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

// Error wrappers
const (
	ErrorCreateConfigFile = "failed to create config file"
	ErrorEncodeConfig     = "failed to encode config"
)

// Config represents the configuration file for dotm.
type Config struct {
	Remote   string
	Repo     string
	Includes []string
	Excludea []string
}

// WriteTomlFile writes a new config file in the toml format to a given path.
func WriteTomlFile(path string, c *Config) error {
	f, err := os.Create(path)
	if err != nil {
		return errors.Wrap(err, ErrorCreateConfigFile)
	}

	w := bufio.NewWriter(f)
	e := toml.NewEncoder(w)
	err = e.Encode(c)
	if err != nil {
		return errors.Wrap(err, ErrorEncodeConfig)
	}

	return nil
}

// NewFromTomlFile reads the file at the given path and decodes it into a new
// config struct.
func NewFromTomlFile(path string) (*Config, error) {
	var config *Config

	data, err := ioutil.ReadFile(path)
	if err != nil {
		// TODO: wrap error
		return nil, err
	}

	_, err = toml.Decode(string(data), config)
	if err != nil {
		// TODO: wrap error
		return nil, err
	}

	return config, nil
}
