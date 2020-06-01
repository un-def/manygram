package profile

import (
	"fmt"
	"io"
	"os"
	"path"
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
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	switch _, err = f.Readdirnames(1); err {
	case nil:
		return nil, fmt.Errorf("%s: already exists", path)
	case io.EOF:
		return &Profile{dir, name, path}, nil
	default:
		return nil, err
	}
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
