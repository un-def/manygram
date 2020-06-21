package desktop

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/suite"
)

const profileName = "foo"
const tryExec = "false"
const exec = "false run foo"

type TestDesktopSuite struct {
	suite.Suite
	dir  string
	path string
}

func (s *TestDesktopSuite) SetupTest() {
	dir, err := ioutil.TempDir("", "test-desktop-*")
	s.Require().NoError(err)
	s.dir = dir
	s.path = path.Join(dir, fmt.Sprintf("telegramdesktop.%s.desktop", profileName))
}

func (s *TestDesktopSuite) TearDownTest() {
	err := os.RemoveAll(s.dir)
	s.Require().NoError(err)
}

func (s *TestDesktopSuite) Create() {
	f, err := os.Create(s.path)
	s.Require().NoError(err)
	f.Close()
}

func (s *TestDesktopSuite) TestExist() {
	exist, err := Exist(s.dir, profileName)
	s.Require().NoError(err)
	s.Require().False(exist)
	s.Create()
	exist, err = Exist(s.dir, profileName)
	s.Require().NoError(err)
	s.Require().True(exist)
}

func (s *TestDesktopSuite) TestCreate() {
	err := Create(s.dir, profileName, tryExec, exec)
	s.Require().NoError(err)
	s.Require().FileExists(s.path)
	contentByte, err := ioutil.ReadFile(s.path)
	s.Require().NoError(err)
	content := string(contentByte)
	name := fmt.Sprintf("Telegram Desktop – %s", profileName)
	s.Require().Contains(content, "Name="+name)
	s.Require().Contains(content, "TryExec="+tryExec)
	s.Require().Contains(content, "Exec="+exec)
}

func (s *TestDesktopSuite) TestRemoveErrNotExist() {
	err := Remove(s.dir, profileName)
	s.Require().Error(err)
	s.Require().True(errors.Is(err, os.ErrNotExist))
}

func (s *TestDesktopSuite) TestRemoveOK() {
	s.Create()
	err := Remove(s.dir, profileName)
	s.Require().NoError(err)
	s.Require().NoFileExists(s.path)
}

func TestDesktopSuiteTest(t *testing.T) {
	suite.Run(t, new(TestDesktopSuite))
}
