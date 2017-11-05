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

// TODO: support env var for CI

// configDir is the dir name where up config is stored relative to HOME.
var configDir = ".up"

// Config is the configuration for the user. The signed Token is the
// source of truth, other values are purely informational.
type Config struct {
	Token string `json:"token"`
	Email string `json:"email"`
	Plan  string `json:"plan"`
}

// Require returns the user config and errors when unauthenticated.
func Require() (*Config, error) {
	var c Config

	if err := c.Load(); err != nil {
		return nil, errors.Wrap(err, "loading config")
	}

	if c.Token == "" {
		return nil, errors.New("user credentials missing, make sure to `up account login` first")
	}

	return &c, nil
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
