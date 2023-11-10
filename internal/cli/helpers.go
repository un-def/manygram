package cli

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/un-def/manygram/internal/config"
	"github.com/un-def/manygram/internal/desktop"
	"github.com/un-def/manygram/internal/profile"
	"github.com/un-def/manygram/internal/xdg"
)

func getDefaultProfileDir(dataDir string) string {
	return path.Join(dataDir, "manygram", "profiles")
}

func getDesktopEntriesDir() string {
	return path.Join(xdg.GetDataHome(), "applications")
}

func printMessage(format string, args ...interface{}) {
	fmt.Fprintf(os.Stdout, format, args...)
	fmt.Fprint(os.Stdout, "\n")
}

func getConfigPath() string {
	return path.Join(xdg.GetConfigHome(), "manygram", "config.toml")
}

func readConfig() (*config.Config, error) {
	configPath := getConfigPath()
	conf, err := config.Read(configPath)
	if err == nil {
		return conf, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return nil, newError(
			"Config %s not found. Run `manygram config create` to create a new one.",
			configPath, err,
		)
	}
	return nil, newError("Failed to read config %s", configPath, err)
}

func readProfile(dir string, name string) (*profile.Profile, error) {
	prof, err := profile.Read(dir, name)
	if err != nil {
		if errors.Is(err, profile.ErrInvalidName) {
			return nil, profileNameError(name)
		}
		if errors.Is(err, profile.ErrNotExist) {
			return nil, newError(
				"Profile '%s' does not exist. Use `manygram create %[1]s` to create a new one.",
				name,
			)
		}
		return nil, newError("Failed to read profile '%s'.", name, err)
	}
	return prof, nil
}

func createDesktopEntry(conf *config.Config, profileName string) error {
	dir := getDesktopEntriesDir()
	exist, err := desktop.Exist(dir, profileName)
	if err != nil {
		return err
	}
	if exist {
		return newError("Desktop entry for profile '%s' already exists.", profileName)
	}
	if conf == nil {
		conf, err = readConfig()
		if err != nil {
			return err
		}
	}
	_, err = readProfile(conf.ProfileDir, profileName)
	if err != nil {
		return err
	}
	if err := desktop.Create(dir, profileName, "manygram", "manygram run "+profileName); err != nil {
		return newError("Failed to create desktop entry for profile '%s'", profileName, err)
	}
	return nil
}

func removeDesktopEntry(profileName string) error {
	err := desktop.Remove(getDesktopEntriesDir(), profileName)
	if err == nil {
		return nil
	} else if errors.Is(err, os.ErrNotExist) {
		return newError("Desktop entry for profile '%s' does not exist.", profileName)
	}
	return newError("Failed to remove desktop entry for profile '%s'.", profileName, err)
}
