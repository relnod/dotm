package config

import (
	"bufio"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	"github.com/relnod/fsa"
	"github.com/relnod/fsa/fsutil"
)

// Errors
var (
	ErrEmptyRemote          = errors.New("empty remote url")
	ErrEmptyPath            = errors.New("empty path")
	ErrProfileAlreadyExists = errors.New("profile already exists")
	ErrProfileNotFound      = errors.New("profile not found")
)

// Error wrappers
const (
	ErrCreateConfigFile = "failed to create config file"
	ErrCreateConfigDir  = "failed to create config directory"
	ErrOpenConfigFile   = "failed to open config file"
	ErrEncodeConfig     = "failed to encode config"
	ErrDecodeConfig     = "failed to decode config"
	ErrInitialzeConfig  = "failed to initialize config"
)

// Config represents the configuration file for dotm.
// A configuration file consists of multiple profiles.
type Config struct {
	FS       fsa.FileSystem      `toml:"-"`
	Profiles map[string]*Profile `toml:"profiles"`
}

// AddProfile adds a new profile to the config. Returns an error if the profile
// already exists, or if one happens during profile initialization.
func (c *Config) AddProfile(name string, p *Profile) error {
	if _, exists := c.Profiles[name]; exists {
		return ErrProfileAlreadyExists
	}

	err := p.Initialize()
	if err != nil {
		return err
	}

	c.Profiles[name] = p

	return nil
}

// FindProfiles tries to find profiles matching a list of profiles.
// If only one name was given and the name is "all", return all profiles.
func (c *Config) FindProfiles(names ...string) (map[string]*Profile, error) {
	if len(names) == 1 && names[0] == "all" {
		return c.Profiles, nil
	}
	profiles := make(map[string]*Profile, 1)
	for _, name := range names {
		p, ok := c.Profiles[name]
		if !ok {
			return nil, ErrProfileNotFound
		}
		profiles[name] = p
	}
	return profiles, nil
}

// Profile defines the configuration for one dotfile forlder.
type Profile struct {
	Remote     string   `toml:"remote"`
	Path       string   `toml:"path"`
	Includes   []string `toml:"includes"`
	Excludes   []string `toml:"excludes"`
	PreUpdate  []string `toml:"pre_update"`
	PostUpdate []string `toml:"post_update"`
}

// Initialize sets the profile up
func (c *Profile) Initialize() error {
	var err error

	c.Remote = os.ExpandEnv(c.Remote)
	c.Path = os.ExpandEnv(c.Path)

	if !filepath.IsAbs(c.Path) {
		c.Path, err = filepath.Abs(c.Path)
		if err != nil {
			return err
		}
	}

	return nil
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
func New(c *Config) (*Config, error) {
	for _, profile := range c.Profiles {
		err := profile.Initialize()
		if err != nil {
			return nil, err
		}
	}
	return c, nil
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
	err := fs.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return errors.Wrap(err, ErrCreateConfigDir)
	}

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
	c := &Config{}

	data, err := fsutil.ReadFile(fs, path)
	if err != nil {
		return nil, errors.Wrap(err, ErrOpenConfigFile)
	}
	_, err = toml.Decode(string(data), c)
	if err != nil {
		return nil, errors.Wrap(err, ErrEncodeConfig)
	}

	c, err = New(c)
	if err != nil {
		return nil, errors.Wrap(err, ErrInitialzeConfig)
	}
	return c, nil
}
