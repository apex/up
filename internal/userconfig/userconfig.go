// Package userconfig provides user machine-level configuration.
package userconfig

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
)

var (
	// configDir is the dir name where up config is stored relative to HOME.
	configDir = ".up"

	// envName is the environment variable which can be used to store
	// Up's configuration, primarily for continuous integration.
	envName = "UP_CONFIG"
)

// Team is the user configuration for a given team.
type Team struct {
	// ID is the team identifier.
	ID string `json:"team"`

	// Email is the user's email.
	Email string `json:"email"`

	// Token is the access token.
	Token string `json:"token"`
}

// IsPersonal returns true if it is a personal team.
func (t *Team) IsPersonal() bool {
	if t.Email == t.ID {
		return true
	}

	return false
}

// Config is the user configuration.
type Config struct {
	// Team is the active team.
	Team string `json:"team"`

	// Teams is the user's active teams.
	Teams map[string]*Team `json:"teams"`
}

// initTeams inits the map.
func (c *Config) initTeams() {
	if c.Teams == nil {
		c.Teams = make(map[string]*Team)
	}
}

// AddTeam adds or replaces the given team.
func (c *Config) AddTeam(t *Team) {
	c.initTeams()
	c.Teams[t.ID] = t
}

// GetTeams returns a list of teams.
func (c *Config) GetTeams() (teams []*Team) {
	for _, t := range c.Teams {
		teams = append(teams, t)
	}
	return
}

// GetTeam returns a team by id or nil
func (c *Config) GetTeam(id string) *Team {
	return c.Teams[id]
}

// GetActiveTeam returns the active team.
func (c *Config) GetActiveTeam() *Team {
	return c.GetTeam(c.Team)
}

// Authenticated returns true if the user has an active team.
func (c *Config) Authenticated() bool {
	return c.GetActiveTeam() != nil
}

// Require requires authentication and returns the active team.
func Require() (*Team, error) {
	var c Config

	if err := c.Load(); err != nil {
		return nil, errors.Wrap(err, "loading config")
	}

	if !c.Authenticated() {
		return nil, errors.New("user credentials missing, make sure to `up account login` first")
	}

	return c.GetActiveTeam(), nil
}

// Alter config, loading and saving after manipulation.
func Alter(fn func(*Config)) error {
	var config Config

	if err := config.Load(); err != nil {
		return errors.Wrap(err, "loading")
	}

	fn(&config)

	if err := config.Save(); err != nil {
		return errors.Wrap(err, "saving")
	}

	return nil
}

// Load the configuration.
func (c *Config) Load() error {
	path, err := c.path()
	if err != nil {
		return errors.Wrap(err, "getting path")
	}

	// env
	if s := os.Getenv(envName); s != "" {
		if err := json.Unmarshal([]byte(s), &c); err != nil {
			return errors.Wrap(err, "unmarshaling")
		}
		return nil
	}

	// file
	b, err := ioutil.ReadFile(path)

	if os.IsNotExist(err) {
		return nil
	}

	if err != nil {
		return errors.Wrap(err, "reading")
	}

	if err := json.Unmarshal(b, c); err != nil {
		return errors.Wrap(err, "unmarshaling")
	}

	return nil
}

// Save the configuration.
func (c *Config) Save() error {
	b, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return errors.Wrap(err, "marshaling")
	}

	path, err := c.path()
	if err != nil {
		return errors.Wrap(err, "getting path")
	}

	if err := ioutil.WriteFile(path, b, 0755); err != nil {
		return errors.Wrap(err, "writing")
	}

	return nil
}

// path returns the path and sets up dir if necessary.
func (c *Config) path() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", errors.Wrap(err, "homedir")
	}

	dir := filepath.Join(home, configDir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", errors.Wrap(err, "mkdir")
	}

	path := filepath.Join(dir, "config.json")
	return path, nil
}
