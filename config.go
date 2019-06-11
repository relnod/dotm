package dotm

import (
	"bufio"
	"errors"
	"fmt"
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
	// Ingoreprefix is a prefix to be ignored when traversing the dotfiles.
	// Both directories and files get ignored if they have the prefix.
	// For directories with the prefix, all sub directories get ignored aswell.
	//
	// The default value is "_".
	IgnorePrefix string `toml:"ignore_prefix" clic:"ignore_prefix"`

	// Profiles is the list of profiles.
	Profiles map[string]*Profile `toml:"profiles" clic:"profile"`
}

// Profile returns the profile with the corresponding name. The returned profile
// has expanded vars.
//
// If no profile with the given name exists an error get returned.
func (c *Config) Profile(name string) (*Profile, error) {
	p, ok := c.Profiles[name]
	if !ok {
		return nil, errors.New("profile does not exist")
	}
	err := p.expandEnv()
	if err != nil {
		return nil, err
	}
	return p, nil
}

// errProfileExists indicates, a profile with the same name already exists
var errProfileExists = errors.New("profile already exists")

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
		return nil, errors.New("profile path already exists")
	}
	if _, err := c.Profile(p.Name); err == nil {
		return nil, errProfileExists
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
		return nil, errProfileExists
	}

	c.Profiles[p.Name] = p
	return p, nil
}

// errOpenConfig indicates that the config file doesn't exist.
var errOpenConfig = errors.New("failed to open config")

// loadConfigWithMetaData tries to load the config file at the given path. Also
// returns the toml metadata.
func loadConfigWithMetaData(path string) (*Config, toml.MetaData, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, toml.MetaData{}, xerrors.Errorf("%v: %w", errOpenConfig, err)
	}
	c := &Config{}
	meta, err := toml.Decode(string(data), c)
	if err != nil {
		return nil, toml.MetaData{}, errors.New("failed to decode config")
	}
	return c, meta, nil
}

// LoadConfig loads the config file.
func LoadConfig() (*Config, error) {
	c, _, err := loadConfigWithMetaData(configPath)
	if err != nil {
		return nil, err
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
		var e *os.PathError
		if xerrors.As(err, &e) {
			return &Config{
				IgnorePrefix: "_", // set default value for the ignored prefix.
				Profiles:     make(map[string]*Profile),
			}, nil
		}
		return nil, err
	}
	return c, nil
}

// Write writes the config file.
func (c *Config) Write() error {
	err := os.MkdirAll(filepath.Dir(configPath), os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create config dir (%s): %v", filepath.Dir(configPath), err)
	}

	f, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("failed to create config file (%s): %v", configPath, err)
	}

	e := toml.NewEncoder(bufio.NewWriter(f))
	err = e.Encode(c)
	if err != nil {
		return errors.New("failed to encode config")
	}

	return nil
}
