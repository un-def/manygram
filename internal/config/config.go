package config

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

// Config ...
type Config struct {
	path       string
	ExecPath   string `toml:"exec-path"`
	ProfileDir string `toml:"profile-dir"`
}

func (c *Config) Write() error {
	dir := filepath.Dir(c.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(c); err != nil {
		return err
	}
	return ioutil.WriteFile(c.path, buf.Bytes(), 0644)
}

// Exist checks whether the config exist
func Exist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else {
		return false, err
	}
}

// New returns new empty config
func New(path string) *Config {
	return &Config{path: path}
}

// Read reads the config from the specified location
func Read(path string) (*Config, error) {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	conf := new(Config)
	md, err := toml.Decode(string(bs), &conf)
	if err != nil {
		return nil, err
	}

	if !md.IsDefined("exec-path") {
		return nil, errors.New("`exec-path` parameter is not defined")
	}
	execPath := strings.TrimSpace(conf.ExecPath)
	if execPath == "" {
		return nil, errors.New("`exec-path` parameter is empty")
	}
	conf.ExecPath = execPath

	if !md.IsDefined("profile-dir") {
		return nil, errors.New("`profile-dir` parameter is not defined")
	}
	profileDir := strings.TrimSpace(conf.ProfileDir)
	if profileDir == "" {
		return nil, errors.New("`profile-dir` parameter is empty")
	}
	conf.ProfileDir = profileDir

	conf.path = path
	return conf, nil
}
