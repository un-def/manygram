package util

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/suite"
)

type TestExistSuite struct {
	suite.Suite
	dir  string
	path string
}

func (s *TestExistSuite) SetupTest() {
	dir, err := ioutil.TempDir("", "test-util-exist-*")
	s.Require().NoError(err)
	s.dir = dir
	s.path = path.Join(dir, "path")
}

func (s *TestExistSuite) TearDownTest() {
	err := os.RemoveAll(s.dir)
	s.Require().NoError(err)
}

func (s *TestExistSuite) TestNotExist() {
	exist, err := Exist(s.path)
	s.Require().NoError(err)
	s.Require().False(exist)
}

func (s *TestExistSuite) TestExistDir() {
	err := os.Mkdir(s.path, 0644)
	s.Require().NoError(err)
	exist, err := Exist(s.path)
	s.Require().NoError(err)
	s.Require().True(exist)
}

func (s *TestExistSuite) TestExistFile() {
	file, err := os.Create(s.path)
	s.Require().NoError(err)
	defer file.Close()
	exist, err := Exist(s.path)
	s.Require().NoError(err)
	s.Require().True(exist)
}

func TestExistSuiteTest(t *testing.T) {
	suite.Run(t, new(TestExistSuite))
}
