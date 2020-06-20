package config

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/suite"
)

type BaseSuite struct {
	suite.Suite
	dir  string
	path string
}

func (s *BaseSuite) SetupTest() {
	dir, err := ioutil.TempDir("", "test-config-read-*")
	s.Require().NoError(err)
	s.dir = dir
	s.path = path.Join(dir, "config.toml")
}

func (s *BaseSuite) TearDownTest() {
	err := os.RemoveAll(s.dir)
	s.Require().NoError(err)
}

func (s *BaseSuite) WriteConfig(content string) {
	err := ioutil.WriteFile(s.path, []byte(content), 0644)
	s.Require().NoError(err)
}

// Read tests

type TestConfigReadSuite struct {
	BaseSuite
}

func (s *TestConfigReadSuite) TestReadNotExist() {
	conf, err := Read(path.Join(s.dir, "conf", "config.toml"))
	s.Require().Error(err)
	s.Require().True(errors.Is(err, os.ErrNotExist), err)
	s.Require().Nil(conf)
}

func (s *TestConfigReadSuite) TestReadIsDirectory() {
	conf, err := Read(s.dir)
	s.Require().Error(err)
	s.Require().Regexp("is a directory", err.Error())
	s.Require().Nil(conf)
}

func (s *TestConfigReadSuite) TestReadMalformed() {
	s.WriteConfig(`bad!format: = 123`)
	conf, err := Read(s.path)
	s.Require().Error(err)
	s.Require().Regexp("bare keys cannot contain", err.Error())
	s.Require().Nil(conf)
}

func (s *TestConfigReadSuite) TestReadNoExecPath() {
	s.WriteConfig(`profile-dir = "/path/to/profiles"`)
	conf, err := Read(s.path)
	s.Require().Error(err)
	s.Require().Regexp("exec-path.*is not defined", err.Error())
	s.Require().Nil(conf)
}

func (s *TestConfigReadSuite) TestReadEmptyExecPath() {
	s.WriteConfig(`
		exec-path = ""
		profile-dir = "/path/to/profiles"
	`)
	conf, err := Read(s.path)
	s.Require().Error(err)
	s.Require().Regexp("exec-path.*is empty", err.Error())
	s.Require().Nil(conf)
}

func (s *TestConfigReadSuite) TestReadNoProfileDir() {
	s.WriteConfig(`exec-path = "/path/to/bin"`)
	conf, err := Read(s.path)
	s.Require().Error(err)
	s.Require().Regexp("profile-dir.*is not defined", err.Error())
	s.Require().Nil(conf)
}

func (s *TestConfigReadSuite) TestReadEmptyProfileDir() {
	s.WriteConfig(`
		exec-path = "/path/to/bin"
		profile-dir = "   "
	`)
	conf, err := Read(s.path)
	s.Require().Error(err)
	s.Require().Regexp("profile-dir.*is empty", err.Error())
	s.Require().Nil(conf)
}

func (s *TestConfigReadSuite) TestReadOK() {
	s.WriteConfig(`
		exec-path = "/path/to/bin"
		profile-dir = "/path/to/profiles"
	`)
	conf, err := Read(s.path)
	s.Require().NoError(err)
	s.Require().Equal(&Config{
		path:       s.path,
		ExecPath:   "/path/to/bin",
		ProfileDir: "/path/to/profiles",
	}, conf)
}

func TestConfigReadSuiteTest(t *testing.T) {
	suite.Run(t, new(TestConfigReadSuite))
}

// Exist tests

type TestConfigExistSuite struct {
	BaseSuite
}

func (s *TestConfigExistSuite) TestNotExist() {
	exist, err := Exist(s.path)
	s.Require().NoError(err)
	s.Require().False(exist)
}

func (s *TestConfigExistSuite) TestExist() {
	s.WriteConfig("")
	exist, err := Exist(s.path)
	s.Require().NoError(err)
	s.Require().True(exist)
}

func TestConfigExistSuiteTest(t *testing.T) {
	suite.Run(t, new(TestConfigExistSuite))
}
