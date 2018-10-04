package config

import (
	"bufio"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	"github.com/relnod/fsa"
	"github.com/relnod/fsa/fsutil"
)

// Errors
var (
	ErrEmptyRemote     = errors.New("empty remote url")
	ErrEmptyPath       = errors.New("empty path")
	ErrProfileNotFound = errors.New("profile not found")
)

// Error wrappers
const (
	ErrCreateConfigFile = "failed to create config file"
	ErrOpenConfigFile   = "failed to open config file"
	ErrEncodeConfig     = "failed to encode config"
	ErrDecodeConfig     = "failed to decode config"
)

// Config represents the configuration file for dotm.
// A configuration file consists of multiple profiles.
type Config struct {
	FS       fsa.FileSystem      `toml:"-"`
	Profiles map[string]*Profile `toml:"profiles"`
}

// FindProfiles tries to find profiles matching a list of profiles.
// If only one name was given and the name is "all", return all profiles.
func (c *Config) FindProfiles(names ...string) ([]*Profile, error) {
	var profiles []*Profile
	if len(names) == 1 && names[0] == "all" {
		for _, p := range c.Profiles {
			profiles = append(profiles, p)
		}
		return profiles, nil
	}
	for _, name := range names {
		p, ok := c.Profiles[name]
		if !ok {
			return nil, ErrProfileNotFound
		}
		profiles = append(profiles, p)
	}
	return profiles, nil
}

// Profile defines the configuration for one dotfile forlder.
type Profile struct {
	Remote   string   `toml:"remote"`
	Path     string   `toml:"path"`
	Includes []string `toml:"includes"`
	Excludes []string `toml:"excludes"`
}

// Validate returns an error if the configuration is invalid.
func (c *Profile) Validate() error {
	if c.Remote == "" {
		// return ErrEmptyRemote
	}
	if c.Path == "" {
		return ErrEmptyPath
	}
	return nil
}

// New takes the given config and intiailizes some values on it.
func New(c *Config) *Config {
	for _, profile := range c.Profiles {
		profile.Remote = os.ExpandEnv(profile.Remote)
		profile.Path = os.ExpandEnv(profile.Path)
	}
	return c
}

// NewConfig returns a new config.
func NewConfig(fs fsa.FileSystem) *Config {
	return &Config{
		FS:       fs,
		Profiles: make(map[string]*Profile, 1),
	}
}

// WriteFile writes a new config file in the toml format to a given path.
func WriteFile(fs fsa.FileSystem, path string, c *Config) error {
	path = os.ExpandEnv(path)
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
	path = os.ExpandEnv(path)
	config := &Config{}

	data, err := fsutil.ReadFile(fs, path)
	if err != nil {
		return nil, errors.Wrap(err, ErrOpenConfigFile)
	}

	_, err = toml.Decode(string(data), config)
	if err != nil {
		return nil, errors.Wrap(err, ErrEncodeConfig)
	}

	return New(config), nil
}
