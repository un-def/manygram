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

type TestNewSuite struct {
	BaseSuite
}

func (s *TestNewSuite) TestOKNotExist() {
	profile, err := New(s.dir, s.name)
	s.Require().NoError(err)
	s.Require().Equal(profile, &Profile{s.dir, s.name, s.path})
	s.Require().DirExists(s.path)
}

func (s *TestNewSuite) TestOKEmpty() {
	s.MakeDir(true)
	profile, err := New(s.dir, s.name)
	s.Require().NoError(err)
	s.Require().Equal(profile, &Profile{s.dir, s.name, s.path})
}

func (s *TestNewSuite) TestErrorNotEmpty() {
	s.MakeDir(false)
	profile, err := New(s.dir, s.name)
	s.Require().Error(err)
	s.Require().Regexp("already exists", err.Error())
	s.Require().Nil(profile)
}

func (s *TestNewSuite) TestErrorIsFile() {
	f, err := os.Create(s.path)
	s.Require().NoError(err)
	f.Close()
	profile, err := New(s.dir, s.name)
	s.Require().Error(err)
	s.Require().Regexp("is a file", err.Error())
	s.Require().Nil(profile)
}

func TestNewSuiteTest(t *testing.T) {
	suite.Run(t, new(TestNewSuite))
}

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
	s.Require().True(errors.Is(err, os.ErrNotExist), err)
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

func TestReaduiteTest(t *testing.T) {
	suite.Run(t, new(TestReadSuite))
}
