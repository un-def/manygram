package tg

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/stretchr/testify/suite"
)

const execName = "fake-telegram-desktop"
const symlinkName = "symlink-telegram-desktop"

type TestExecutableSuite struct {
	suite.Suite
	dir      string
	execPath string
}

func (s *TestExecutableSuite) SetupTest() {
	dir, err := ioutil.TempDir("", "test-tg-executable-*")
	s.Require().NoError(err)
	s.dir = dir
	s.execPath = path.Join(dir, execName)
}

func (s *TestExecutableSuite) TearDownTest() {
	err := os.RemoveAll(s.dir)
	s.Require().NoError(err)
}

func (s *TestExecutableSuite) ErrorIs(err error, target error) {
	s.Require().True(errors.Is(err, target), err)
}

func (s *TestExecutableSuite) CreateFile(path string, isExecutable bool) {
	var perm os.FileMode
	if isExecutable {
		perm = 0777
	} else {
		perm = 0666
	}
	f, err := os.OpenFile(path, os.O_CREATE, perm)
	s.Require().NoError(err)
	f.Close()
}

func (s *TestExecutableSuite) CreateSymlink(src string, dest string) {
	err := os.Symlink(src, dest)
	s.Require().NoError(err)
}

func (s *TestExecutableSuite) TestErrPathNotExist() {
	tg, err := Executable(s.execPath, nil)
	s.Require().Error(err)
	s.ErrorIs(err, os.ErrNotExist)
	s.Require().Nil(tg)
}

func (s *TestExecutableSuite) TestErrPathIsDir() {
	tg, err := Executable(s.dir, nil)
	s.Require().Error(err)
	s.ErrorIs(err, os.ErrPermission)
	s.Require().Nil(tg)
}

func (s *TestExecutableSuite) TestErrPathIsNotExecutable() {
	s.CreateFile(s.execPath, false)
	tg, err := Executable(s.execPath, nil)
	s.Require().Error(err)
	s.ErrorIs(err, os.ErrPermission)
	s.Require().Nil(tg)
}

func (s *TestExecutableSuite) TestErrNameNotExist() {
	tg, err := Executable(execName, nil)
	s.Require().Error(err)
	s.ErrorIs(err, exec.ErrNotFound)
	s.Require().Nil(tg)
}

func (s *TestExecutableSuite) TestOKPath() {
	s.CreateFile(s.execPath, true)
	tg, err := Executable(s.execPath, nil)
	s.Require().NoError(err)
	s.Require().Equal(&TelegramDesktop{
		s.execPath,
		s.execPath,
		s.execPath,
		nil,
	}, tg)
}

func (s *TestExecutableSuite) TestOKPathSymlink() {
	s.CreateFile(s.execPath, true)
	symlinkPath := path.Join(s.dir, symlinkName)
	s.CreateSymlink(s.execPath, symlinkPath)
	tg, err := Executable(symlinkPath, nil)
	s.Require().NoError(err)
	s.Require().Equal(&TelegramDesktop{
		symlinkPath,
		symlinkPath,
		s.execPath,
		nil,
	}, tg)
}

func (s *TestExecutableSuite) TestOKName() {
	s.CreateFile(s.execPath, true)
	origPATH := os.Getenv("PATH")
	defer func() { os.Setenv("PATH", origPATH) }()
	os.Setenv("PATH", fmt.Sprintf("%s:%s", s.dir, origPATH))
	tg, err := Executable(execName, nil)
	s.Require().NoError(err)
	s.Require().Equal(&TelegramDesktop{
		execName,
		s.execPath,
		s.execPath,
		nil,
	}, tg)
}

func (s *TestExecutableSuite) TestOKWithArgs() {
	s.CreateFile(s.execPath, true)
	tg, err := Executable(s.execPath, []string{"-extra1", "-extra2"})
	s.Require().NoError(err)
	s.Require().Equal(&TelegramDesktop{
		s.execPath,
		s.execPath,
		s.execPath,
		[]string{"-extra1", "-extra2"},
	}, tg)
}

func (s *TestExecutableSuite) TestIsSnapFalse() {
	s.CreateFile(s.execPath, true)
	symlinkPath := path.Join(s.dir, symlinkName)
	s.CreateSymlink(s.execPath, symlinkPath)
	tg, err := Executable(symlinkPath, nil)
	s.Require().NoError(err)
	s.Require().False(tg.IsSnap())
}

func (s *TestExecutableSuite) TestIsSnapTrue() {
	snapPath := path.Join(s.dir, "snap")
	s.CreateFile(snapPath, true)
	symlinkPath := path.Join(s.dir, symlinkName)
	s.CreateSymlink(snapPath, symlinkPath)
	tg, err := Executable(symlinkPath, nil)
	s.Require().NoError(err)
	s.Require().True(tg.IsSnap())
}

func TestExecutableSuiteTest(t *testing.T) {
	suite.Run(t, new(TestExecutableSuite))
}

type TestFlatpakSuite struct {
	suite.Suite
	dir         string
	flatpakPath string
	origPATH    string
}

func (s *TestFlatpakSuite) SetupTest() {
	dir, err := ioutil.TempDir("", "test-bin-*")
	s.Require().NoError(err)
	s.dir = dir
	s.flatpakPath = path.Join(s.dir, "flatpak")
	s.origPATH = os.Getenv("PATH")
	os.Setenv("PATH", dir)
}

func (s *TestFlatpakSuite) TearDownTest() {
	os.Setenv("PATH", s.origPATH)
	err := os.RemoveAll(s.dir)
	s.Require().NoError(err)
}

func (s *TestFlatpakSuite) CreateFlatpakExecutable(content ...string) {
	f, err := os.OpenFile(s.flatpakPath, os.O_CREATE|os.O_WRONLY, 0777)
	s.Require().NoError(err)
	defer f.Close()
	w := bufio.NewWriter(f)
	_, err = f.WriteString("#!/bin/sh\n")
	s.Require().NoError(err)
	for _, line := range content {
		_, err = w.WriteString(line)
		s.Require().NoError(err)
		_, err = w.WriteRune('\n')
		s.Require().NoError(err)
	}
	w.Flush()
}

func (s *TestFlatpakSuite) TestErrFlatpakNotFound() {
	tg, err := Flatpak()
	s.Require().ErrorContains(err, "flatpak executable not found")
	s.Require().Nil(tg)
}

func (s *TestFlatpakSuite) TestErrFlatpakAppNotFound() {
	s.CreateFlatpakExecutable(`exit 1`)
	tg, err := Flatpak()
	s.Require().ErrorContains(err, "org.telegram.desktop flatpak app not found")
	s.Require().Nil(tg)
}

func (s *TestFlatpakSuite) TestOKUser() {
	s.CreateFlatpakExecutable(
		`while [ $# -gt 0 ]; do`,
		`  if [ "$1" = '--user' ]; then exit 0; fi`,
		`  shift`,
		`done`,
		`exit 1`,
	)
	tg, err := Flatpak()
	s.Require().NoError(err)
	s.Require().Equal(&TelegramDesktop{
		"flatpak",
		s.flatpakPath,
		s.flatpakPath,
		[]string{"run", "--user", "org.telegram.desktop"},
	}, tg)
}

func (s *TestFlatpakSuite) TestOKSystem() {
	s.CreateFlatpakExecutable(
		`while [ $# -gt 0 ]; do`,
		`  if [ "$1" = '--user' ]; then exit 1; fi`,
		`  shift`,
		`done`,
		`exit 0`,
	)
	tg, err := Flatpak()
	s.Require().NoError(err)
	s.Require().Equal(&TelegramDesktop{
		"flatpak",
		s.flatpakPath,
		s.flatpakPath,
		[]string{"run", "org.telegram.desktop"},
	}, tg)
}

func TestFlatpakSuiteTest(t *testing.T) {
	suite.Run(t, new(TestFlatpakSuite))
}

type BaseFakeHOMESuite struct {
	suite.Suite
	dir      string
	origHOME string
}

func (s *BaseFakeHOMESuite) SetupTest() {
	dir, err := ioutil.TempDir("", "fake-home-dir-*")
	s.Require().NoError(err)
	s.dir = dir
	s.origHOME = os.Getenv("HOME")
	os.Setenv("HOME", dir)
}

func (s *BaseFakeHOMESuite) TearDownTest() {
	os.Setenv("HOME", s.origHOME)
	err := os.RemoveAll(s.dir)
	s.Require().NoError(err)
}

type TestGetFlatpakDataHomeSuite struct {
	BaseFakeHOMESuite
}

func (s *TestGetFlatpakDataHomeSuite) TestErr() {
	path, err := GetFlatpakDataHome()
	s.Require().Error(err)
	s.Require().Equal("", path)
}

func (s *TestGetFlatpakDataHomeSuite) TestOK() {
	dataHome := path.Join(s.dir, ".var/app/org.telegram.desktop/data")
	os.MkdirAll(dataHome, 0777)
	path, err := GetFlatpakDataHome()
	s.Require().NoError(err)
	s.Require().Equal(dataHome, path)
}

func TestGetFlatpakDataHomeSuiteTest(t *testing.T) {
	suite.Run(t, new(TestGetFlatpakDataHomeSuite))
}

type TestGetSnapDataHomeSuite struct {
	BaseFakeHOMESuite
}

func (s *TestGetSnapDataHomeSuite) TestErr() {
	path, err := GetSnapDataHome()
	s.Require().Error(err)
	s.Require().Equal("", path)
}

func (s *TestGetSnapDataHomeSuite) TestOK() {
	dataHome := path.Join(s.dir, "snap/telegram-desktop/current/.local/share")
	os.MkdirAll(dataHome, 0777)
	path, err := GetSnapDataHome()
	s.Require().NoError(err)
	s.Require().Equal(dataHome, path)
}

func TestGetSnapDataHomeSuiteTest(t *testing.T) {
	suite.Run(t, new(TestGetSnapDataHomeSuite))
}
