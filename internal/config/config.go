package config

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

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
		return nil, errors.New("exec-path is not defined")
	}
	conf.path = path
	return conf, nil
}

// Default returns the default config
func Default(path string) *Config {
	return &Config{
		path:       path,
		ExecPath:   "telegram-desktop",
		ProfileDir: "",
	}
}
