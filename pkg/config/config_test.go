package config_test

import (
	"testing"

	"github.com/relnod/dotm/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestConfigFindProfiles(t *testing.T) {
	profile1 := &config.Profile{Path: "a"}
	profile2 := &config.Profile{Path: "b"}
	profile3 := &config.Profile{Path: "c"}

	c := &config.Config{
		Profiles: map[string]*config.Profile{
			"test1": profile1,
			"test2": profile2,
			"test3": profile3,
		},
	}

	var profiles []*config.Profile
	var err error

	// Find profile test1
	profiles, err = c.FindProfiles("test1")
	assert.NoError(t, err)
	assert.EqualValues(t, []*config.Profile{profile1}, profiles)

	// Find profiles test2 and test3
	profiles, err = c.FindProfiles("test2", "test3")
	assert.NoError(t, err)
	assert.EqualValues(t, []*config.Profile{profile2, profile3}, profiles)

	// Find all profiles
	profiles, err = c.FindProfiles("all")
	assert.NoError(t, err)
	assert.EqualValues(t, []*config.Profile{profile1, profile2, profile3}, profiles)

	// Error, when profile was not found
	profiles, err = c.FindProfiles("blablub")
	assert.Equal(t, config.ErrProfileNotFound, err)
	assert.Nil(t, profiles)
}
