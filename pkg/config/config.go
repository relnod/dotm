package config

import (
	"bufio"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	"github.com/relnod/fsa"
	"github.com/relnod/fsa/fsutil"

	"github.com/relnod/dotm/pkg/profile"
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
	FS       fsa.FileSystem              `toml:"-"`
	Profiles map[string]*profile.Profile `toml:"profiles"`
}

// NewConfig returns a new config.
func NewConfig(fs fsa.FileSystem) *Config {
	return &Config{
		FS:       fs,
		Profiles: make(map[string]*profile.Profile, 1),
	}
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

	c.FS = fs
	err = c.InitializeProfiles()
	if err != nil {
		return nil, errors.Wrap(err, ErrInitialzeConfig)
	}
	return c, nil
}

// InitializeProfiles initializes the profiles.
func (c *Config) InitializeProfiles() error {
	for name, p := range c.Profiles {
		p.SetFS(c.FS)
		err := p.ExpandEnvs(name)
		if err != nil {
			return err
		}
	}
	return nil
}

// AddProfile adds a new profile to the config. Returns an error if the profile
// already exists, or if one happens during profile initialization.
func (c *Config) AddProfile(name string, p *profile.Profile) error {
	if _, exists := c.Profiles[name]; exists {
		return ErrProfileAlreadyExists
	}

	p.SetFS(c.FS)
	err := p.ExpandEnvs(name)
	if err != nil {
		return err
	}

	c.Profiles[name] = p

	return nil
}

// FindProfile tries to find a profile with a given name.
func (c *Config) FindProfile(name string) (*profile.Profile, error) {
	p, ok := c.Profiles[name]
	if !ok {
		return nil, ErrProfileNotFound
	}
	return p, nil
}

// FindProfiles tries to find profiles matching a list of profiles.
// If only one name was given and the name is "all", return all profiles.
func (c *Config) FindProfiles(names ...string) (map[string]*profile.Profile, error) {
	if len(names) == 1 && names[0] == "all" {
		return c.Profiles, nil
	}
	profiles := make(map[string]*profile.Profile, 1)
	for _, name := range names {
		p, err := c.FindProfile(name)
		if err != nil {
			return nil, err
		}
		profiles[name] = p
	}
	return profiles, nil
}

// WriteFile writes a new config file in the toml format to a given path.
func (c *Config) WriteFile(path string) error {
	path = os.ExpandEnv(path)
	err := c.FS.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return errors.Wrap(err, ErrCreateConfigDir)
	}

	f, err := c.FS.Create(path)
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
