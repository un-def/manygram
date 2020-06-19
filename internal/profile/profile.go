package profile

import (
	"fmt"
	"os"
	"path"
	"strings"
)

// Profile type
type Profile struct {
	Dir  string
	Name string
	Path string
}

// New creates a new profile directory
func New(dir string, name string) (*Profile, error) {
	path := path.Join(dir, name)
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return nil, err
		}
		return &Profile{dir, name, path}, nil
	} else if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%s: is a file", path)
	}
	return nil, fmt.Errorf("%s: already exists", path)
}

// Read checks if the profile directory exists
func Read(dir string, name string) (*Profile, error) {
	path := path.Join(dir, name)
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%s: is a file", path)
	}
	return &Profile{dir, name, path}, nil
}

// IsProfileDirExist checks whether profile directory exists
func IsProfileDirExist(path string) (bool, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return false, fmt.Errorf("%s: empty path", path)
	}
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if !info.IsDir() {
		return false, fmt.Errorf("%s: is a file", path)
	}
	return true, nil
}
