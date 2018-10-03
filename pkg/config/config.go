package config

import (
	"bufio"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	"github.com/relnod/fsa"
	"github.com/relnod/fsa/fsutil"
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
	FS       fsa.FileSystem
}

// Validate returns an error if the configuration is invalid.
func (c *Config) Validate() error {
	if c.Remote == "" {
		// return ErrEmptyRemote
	}
	if c.Path == "" {
		return ErrEmptyPath
	}
	return nil
}

// WriteFile writes a new config file in the toml format to a given path.
func WriteFile(fs fsa.FileSystem, path string, c *Config) error {
	f, err := fs.Create(path)
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

// NewFromFile reads the file at the given path and decodes it into a new
// config struct.
func NewFromFile(fs fsa.FileSystem, path string) (*Config, error) {
	config := &Config{}

	data, err := fsutil.ReadFile(fs, path)
	if err != nil {
		return nil, errors.Wrap(err, ErrOpenConfigFile)
	}

	_, err = toml.Decode(string(data), config)
	if err != nil {
		return nil, errors.Wrap(err, ErrEncodeConfig)
	}

	return config, nil
}
