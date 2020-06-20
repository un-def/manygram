package profile

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
	name string
	path string
}

func (s *BaseSuite) SetupTest() {
	dir, err := ioutil.TempDir("", "test-profile-*")
	s.Require().NoError(err)
	s.dir = dir
	s.name = "username"
	s.path = path.Join(dir, s.name)
}

func (s *BaseSuite) TearDownTest() {
	err := os.RemoveAll(s.dir)
	s.Require().NoError(err)
}

func (s *BaseSuite) MakeDir(empty bool) {
	err := os.Mkdir(s.path, 0755)
	s.Require().NoError(err)
	if empty {
		return
	}
	f, err := os.Create(path.Join(s.path, "some-file"))
	s.Require().NoError(err)
	f.Close()
}

// Create tests

type TestCreateSuite struct {
	BaseSuite
}

func (s *TestCreateSuite) TestOKNotExist() {
	profile, err := Create(s.dir, s.name)
	s.Require().NoError(err)
	s.Require().Equal(profile, &Profile{s.dir, s.name, s.path})
	s.Require().DirExists(s.path)
}

func (s *TestCreateSuite) TestErrorExistEmpty() {
	s.MakeDir(true)
	profile, err := Create(s.dir, s.name)
	s.Require().Error(err)
	s.Require().Regexp("already exists", err.Error())
	s.Require().Nil(profile)
}

func (s *TestCreateSuite) TestErrorExistNotEmpty() {
	s.MakeDir(false)
	profile, err := Create(s.dir, s.name)
	s.Require().Error(err)
	s.Require().Regexp("already exists", err.Error())
	s.Require().Nil(profile)
}

func (s *TestCreateSuite) TestErrorIsFile() {
	f, err := os.Create(s.path)
	s.Require().NoError(err)
	f.Close()
	profile, err := Create(s.dir, s.name)
	s.Require().Error(err)
	s.Require().Regexp("is a file", err.Error())
	s.Require().Nil(profile)
}

func (s *TestCreateSuite) TestErrInvalidName() {
	profile, err := Create(s.dir, "foo/bar")
	s.Require().Error(err)
	s.Require().Regexp("invalid profile name", err.Error())
	s.Require().Nil(profile)
}

func TestCreateSuiteTest(t *testing.T) {
	suite.Run(t, new(TestCreateSuite))
}

// Read tests

type TestReadSuite struct {
	BaseSuite
}

func (s *TestReadSuite) TestOKEmpty() {
	s.MakeDir(true)
	profile, err := Read(s.dir, s.name)
	s.Require().NoError(err)
	s.Require().Equal(profile, &Profile{s.dir, s.name, s.path})
}

func (s *TestReadSuite) TestOKNotEmpty() {
	s.MakeDir(false)
	profile, err := Read(s.dir, s.name)
	s.Require().NoError(err)
	s.Require().Equal(profile, &Profile{s.dir, s.name, s.path})
}

func (s *TestReadSuite) TestErrorNotExist() {
	profile, err := Read(s.dir, s.name)
	s.Require().Error(err)
	s.Require().True(errors.Is(err, ErrNotExist), err)
	s.Require().Nil(profile)
}

func (s *TestReadSuite) TestErrorIsFile() {
	f, err := os.Create(s.path)
	s.Require().NoError(err)
	f.Close()
	profile, err := Read(s.dir, s.name)
	s.Require().Error(err)
	s.Require().Regexp("is a file", err.Error())
	s.Require().Nil(profile)
}

func (s *TestReadSuite) TestErrorInvalidName() {
	profile, err := Read(s.dir, "1foobar")
	s.Require().Error(err)
	s.Require().Regexp("invalid profile name", err.Error())
	s.Require().Nil(profile)
}

func TestReadSuiteTest(t *testing.T) {
	suite.Run(t, new(TestReadSuite))
}

// Remove tests

type TestRemoveSuite struct {
	BaseSuite
}

func (s *TestRemoveSuite) TestErrorNotExist() {
	err := Remove(s.dir, "non_existent")
	s.Require().Error(err)
	s.Require().True(errors.Is(err, ErrNotExist), err)
}

func (s *TestRemoveSuite) TestOK() {
	s.MakeDir(false)
	err := Remove(s.dir, s.name)
	s.Require().NoError(err)
	_, err = os.Stat(s.path)
	s.Require().Error(err)
	s.Require().True(errors.Is(err, ErrNotExist), err)
}

func (s *TestRemoveSuite) TestErrorInvalidName() {
	err := Remove(s.dir, "")
	s.Require().Error(err)
	s.Require().Regexp("invalid profile name", err.Error())
}

func TestRemoveSuiteTest(t *testing.T) {
	suite.Run(t, new(TestRemoveSuite))
}

// IsProfileDirExist tests

type TestIsProfileDirExistSuite struct {
	suite.Suite
	dir string
}

func (s *TestIsProfileDirExistSuite) SetupTest() {
	dir, err := ioutil.TempDir("", "test-profile-dir-*")
	s.Require().NoError(err)
	s.dir = dir
}

func (s *TestIsProfileDirExistSuite) TearDownTest() {
	err := os.RemoveAll(s.dir)
	s.Require().NoError(err)
}

func (s *TestIsProfileDirExistSuite) TestOKNotExist() {
	exist, err := IsProfileDirExist(path.Join(s.dir, "profiles"))
	s.Require().NoError(err)
	s.Require().False(exist)
}

func (s *TestIsProfileDirExistSuite) TestOKExist() {
	exist, err := IsProfileDirExist(s.dir)
	s.Require().NoError(err)
	s.Require().True(exist)
}

func (s *TestIsProfileDirExistSuite) TestErrIsFile() {
	path := path.Join(s.dir, "profiles")
	f, err := os.Create(path)
	s.Require().NoError(err)
	f.Close()
	exist, err := IsProfileDirExist(path)
	s.Require().Error(err)
	s.Require().Regexp("is a file", err.Error())
	s.Require().False(exist)
}

func TestIsProfileDirExistSuiteTest(t *testing.T) {
	suite.Run(t, new(TestIsProfileDirExistSuite))
}
