package dotm

import (
	"bufio"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"golang.org/x/xerrors"
)

// configPath is the path to the configuration file
var configPath = os.ExpandEnv("$HOME/.config/dotm/config.toml")

// Config represents the configuration file for dotm.
// A configuration file consists of multiple profiles.
type Config struct {
	Profiles map[string]*Profile `toml:"profiles"`
}

// ErrProfileNotExists indicates, that the dotfile profile was not declared in
// the config file.
var ErrProfileNotExists = errors.New("profile does not exist")

// Profile returns the profile with the corresponding name. If no profile with
// the given name exists an error get returned.
func (c *Config) Profile(name string) (*Profile, error) {
	p, ok := c.Profiles[name]
	if !ok {
		return nil, ErrProfileNotExists
	}
	return p, nil
}

// ErrOpenConfig indicates that the config file doesn't exist.
var ErrOpenConfig = xerrors.Errorf("failed to open config")

// ErrDecodeConfig indicates that there is a syntax error in the config file.
var ErrDecodeConfig = xerrors.Errorf("failed to decode config")

// LoadConfig loads the config file.
func LoadConfig() (*Config, error) {
	c := &Config{}

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, ErrOpenConfig
	}
	_, err = toml.Decode(string(data), c)
	if err != nil {
		return nil, ErrDecodeConfig
	}

	for name, p := range c.Profiles {
		p.Name = name
		err := p.expandEnv()
		if err != nil {
			return nil, err
		}
	}
	return c, nil
}

// LoadOrCreateConfig tries to load the config file. If it doesn't not exist, it
// returns a new config.
func LoadOrCreateConfig() (*Config, error) {
	c, err := LoadConfig()
	if err != nil {
		if err == ErrOpenConfig {
			return &Config{make(map[string]*Profile)}, nil
		}
		return nil, err
	}
	return c, nil
}

// ErrCreateConfigDir indicates a failure during the creation of the config dir.
var ErrCreateConfigDir = xerrors.Errorf("failed to create config dir (%s)", filepath.Dir(configPath))

// ErrCreateConfigFile indicates a failure during the creation file.
var ErrCreateConfigFile = xerrors.Errorf("failed to create config file (%s)", configPath)

// ErrEncodeConfig indicates a failure during the encoding of the config file.
var ErrEncodeConfig = xerrors.Errorf("failed to encode config")

// Write writes the config file.
func (c *Config) Write() error {
	err := os.MkdirAll(filepath.Dir(configPath), os.ModePerm)
	if err != nil {
		return ErrCreateConfigDir
	}

	f, err := os.Create(configPath)
	if err != nil {
		return ErrCreateConfigFile
	}

	e := toml.NewEncoder(bufio.NewWriter(f))
	err = e.Encode(c)
	if err != nil {
		return ErrEncodeConfig
	}

	return nil
}
