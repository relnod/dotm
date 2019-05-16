package dotm

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// configPath is the path to the configuration file
var configPath = os.ExpandEnv("$HOME/.config/dotm/config.toml")

// Config represents the configuration file for dotm.
// A configuration file consists of multiple profiles.
type Config struct {
	// Ingoreprefix is a prefix to be ignored when traversing the dotfiles.
	// Both directories and files get ignored if they have the prefix.
	// For directories with the prefix, all sub directories get ignored aswell.
	//
	// The default value is "_".
	IgnorePrefix string `toml:"ignore_prefix"`

	// Profiles is the list of profiles.
	Profiles map[string]*Profile `toml:"profiles"`
}

// ErrProfileNotExists indicates, that the dotfile profile was not declared in
// the config file.
var ErrProfileNotExists = errors.New("profile does not exist")

// Profile returns the profile with the corresponding name. The returned profile
// has expanded vars.
//
// If no profile with the given name exists an error get returned.
func (c *Config) Profile(name string) (*Profile, error) {
	p, ok := c.Profiles[name]
	if !ok {
		return nil, ErrProfileNotExists
	}
	err := p.expandEnv()
	if err != nil {
		return nil, err
	}
	return p, nil
}

// ErrProfilePathExists indicates, that the profile path already exists.
var ErrProfilePathExists = errors.New("profile path already exists")

// ErrProfileExists indicates, a profile with the same name already exists
var ErrProfileExists = errors.New("profile already exists")

// AddProfile adds a new profile. If a profile with the same name exists, it
// returns an error.
//
// It sanitizes and expands vars.
//
// If the profile path already exists, it returns an error.
// If a profile with the same name exists, it returns an error.
func (c *Config) AddProfile(p *Profile) (*Profile, error) {
	p.sanitize()
	if err := p.expandEnv(); err != nil {
		return nil, err
	}

	if _, err := os.Stat(p.expandedPath); err == nil {
		return nil, ErrProfilePathExists
	}
	if _, err := c.Profile(p.Name); err == nil {
		return nil, ErrProfileExists
	}

	c.Profiles[p.Name] = p
	return p, nil
}

// AddProfileFromExistingPath adds a new profile.
//
// It sanitizes and expands vars.
//
// If a profile with the same name exists, it returns an error.
func (c *Config) AddProfileFromExistingPath(p *Profile) (*Profile, error) {
	p.sanitize()
	if err := p.expandEnv(); err != nil {
		return nil, err
	}

	if _, err := c.Profile(p.Name); err == nil {
		return nil, ErrProfileExists
	}

	c.Profiles[p.Name] = p
	return p, nil
}

// ErrOpenConfig indicates that the config file doesn't exist.
var ErrOpenConfig = errors.New("failed to open config")

// ErrDecodeConfig indicates that there is a syntax error in the config file.
var ErrDecodeConfig = errors.New("failed to decode config")

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
			return &Config{
				IgnorePrefix: "_", // set default value for the ignored prefix.
				Profiles:     make(map[string]*Profile),
			}, nil
		}
		return nil, err
	}
	return c, nil
}

// ErrCreateConfigDir indicates a failure during the creation of the config dir.
var ErrCreateConfigDir = fmt.Errorf("failed to create config dir (%s)", filepath.Dir(configPath))

// ErrCreateConfigFile indicates a failure during the creation file.
var ErrCreateConfigFile = fmt.Errorf("failed to create config file (%s)", configPath)

// ErrEncodeConfig indicates a failure during the encoding of the config file.
var ErrEncodeConfig = errors.New("failed to encode config")

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
