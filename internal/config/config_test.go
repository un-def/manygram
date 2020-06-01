package config

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadNotExist(t *testing.T) {
	dir, err := ioutil.TempDir("", "test-config-*")
	require.NoError(t, err)
	defer os.RemoveAll(dir)
	conf, err := Read(path.Join(dir, "conf", "config.toml"))
	require.Error(t, err)
	require.True(t, errors.Is(err, os.ErrNotExist), err)
	require.Nil(t, conf)
}

func TestReadIsDirectory(t *testing.T) {
	dir, err := ioutil.TempDir("", "test-config-*")
	require.NoError(t, err)
	defer os.RemoveAll(dir)
	conf, err := Read(dir)
	require.Error(t, err)
	require.Regexp(t, "is a directory", err.Error())
	require.Nil(t, conf)
}

func TestReadMalformed(t *testing.T) {
	content := []byte(`bad!format: = 123`)
	tmpfile, err := ioutil.TempFile("", "test-config-*.toml")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()
	_, err = tmpfile.Write(content)
	require.NoError(t, err)
	conf, err := Read(tmpfile.Name())
	require.Error(t, err)
	require.Regexp(t, "bare keys cannot contain", err.Error())
	require.Nil(t, conf)
}

func TestReadOK(t *testing.T) {
	content := []byte(`
		bin-path = "/path/to/bin"
		profile-dir = "/path/to/profiles"
	`)
	tmpfile, err := ioutil.TempFile("", "test-config-*.toml")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()
	_, err = tmpfile.Write(content)
	require.NoError(t, err)
	path := tmpfile.Name()
	conf, err := Read(path)
	require.NoError(t, err)
	require.Equal(t, &Config{
		path:       path,
		BinPath:    "/path/to/bin",
		ProfileDir: "/path/to/profiles",
	}, conf)
}
