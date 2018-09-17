package config

import (
	"bufio"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

// Errors
var (
	ErrEmptyRemote = errors.New("empty remote url")
	ErrEmptyPath   = errors.New("empty path")
)

// Error wrappers
const (
	ErrCreateConfigFile = "failed to create config file"
	ErrOpenConfigFile   = "failed to open config file"
	ErrEncodeConfig     = "failed to encode config"
	ErrDecodeConfig     = "failed to decode config"
)

// Config represents the configuration file for dotm.
type Config struct {
	Remote   string
	Path     string
	Includes []string
	Excludea []string
}

// Validate returns an error if the configuration is invalid.
func (c *Config) Validate() error {
	if c.Remote == "" {
		return ErrEmptyRemote
	}
	if c.Path == "" {
		return ErrEmptyPath
	}
	return nil
}

// WriteTomlFile writes a new config file in the toml format to a given path.
func WriteTomlFile(path string, c *Config) error {
	f, err := os.Create(path)
	if err != nil {
		return errors.Wrap(err, ErrCreateConfigFile)
	}

	w := bufio.NewWriter(f)
	e := toml.NewEncoder(w)
	err = e.Encode(c)
	if err != nil {
		return errors.Wrap(err, ErrEncodeConfig)
	}

	return nil
}

// NewFromTomlFile reads the file at the given path and decodes it into a new
// config struct.
func NewFromTomlFile(path string) (*Config, error) {
	config := &Config{}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, ErrOpenConfigFile)
	}

	_, err = toml.Decode(string(data), config)
	if err != nil {
		return nil, errors.Wrap(err, ErrEncodeConfig)
	}

	return config, nil
}
