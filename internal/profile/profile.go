package profile

import (
	"errors"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
)

// Profile type
type Profile struct {
	Dir  string
	Name string
	Path string
}

// ErrAlreadyExists is returned by the New() function when the profile directory exists
var ErrAlreadyExists = errors.New("already exists")

// ErrNotExist is returned by the Read() function when the profile directory does not exist
var ErrNotExist = os.ErrNotExist

// ErrInvalidName indicates that the profile name does not meet requirements
var ErrInvalidName = errors.New("invalid profile name")

var profileNameRegexp = regexp.MustCompile("^[A-Za-z][A-Za-z0-9_]*$")

// New creates a new profile directory
func New(dir string, name string) (*Profile, error) {
	if !IsValidName(name) {
		return nil, ErrInvalidName
	}
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
	return nil, fmt.Errorf("%s: %w", path, ErrAlreadyExists)
}

// Read checks if the profile directory exists
func Read(dir string, name string) (*Profile, error) {
	if !IsValidName(name) {
		return nil, ErrInvalidName
	}
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

// IsValidName checks whether the profile name meets requirements
func IsValidName(name string) bool {
	return profileNameRegexp.MatchString(name)
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
