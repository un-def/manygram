package tg

import (
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

type TestTelegramSuite struct {
	suite.Suite
	dir      string
	execPath string
}

func (s *TestTelegramSuite) SetupTest() {
	dir, err := ioutil.TempDir("", "test-tg-executable-*")
	s.Require().NoError(err)
	s.dir = dir
	s.execPath = path.Join(dir, execName)
}

func (s *TestTelegramSuite) TearDownTest() {
	err := os.RemoveAll(s.dir)
	s.Require().NoError(err)
}

func (s *TestTelegramSuite) ErrorIs(err error, target error) {
	s.Require().True(errors.Is(err, target), err)
}

func (s *TestTelegramSuite) CreateFile(path string, isExecutable bool) {
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

func (s *TestTelegramSuite) CreateSymlink(src string, dest string) {
	err := os.Symlink(src, dest)
	s.Require().NoError(err)
}

func (s *TestTelegramSuite) TestErrPathNotExist() {
	tg, err := Executable(s.execPath)
	s.Require().Error(err)
	s.ErrorIs(err, os.ErrNotExist)
	s.Require().Nil(tg)
}

func (s *TestTelegramSuite) TestErrPathIsDir() {
	tg, err := Executable(s.dir)
	s.Require().Error(err)
	s.ErrorIs(err, os.ErrPermission)
	s.Require().Nil(tg)
}

func (s *TestTelegramSuite) TestErrPathIsNotExecutable() {
	s.CreateFile(s.execPath, false)
	tg, err := Executable(s.execPath)
	s.Require().Error(err)
	s.ErrorIs(err, os.ErrPermission)
	s.Require().Nil(tg)
}

func (s *TestTelegramSuite) TestErrNameNotExist() {
	tg, err := Executable(execName)
	s.Require().Error(err)
	s.ErrorIs(err, exec.ErrNotFound)
	s.Require().Nil(tg)
}

func (s *TestTelegramSuite) TestOKPath() {
	s.CreateFile(s.execPath, true)
	tg, err := Executable(s.execPath)
	s.Require().NoError(err)
	s.Require().Equal(&TelegramDesktop{
		s.execPath,
		s.execPath,
		s.execPath,
	}, tg)
}

func (s *TestTelegramSuite) TestOKPathSymlink() {
	s.CreateFile(s.execPath, true)
	symlinkPath := path.Join(s.dir, symlinkName)
	s.CreateSymlink(s.execPath, symlinkPath)
	tg, err := Executable(symlinkPath)
	s.Require().NoError(err)
	s.Require().Equal(&TelegramDesktop{
		symlinkPath,
		symlinkPath,
		s.execPath,
	}, tg)
}

func (s *TestTelegramSuite) TestOKName() {
	s.CreateFile(s.execPath, true)
	origPATH := os.Getenv("PATH")
	defer func() { os.Setenv("PATH", origPATH) }()
	os.Setenv("PATH", fmt.Sprintf("%s:%s", s.dir, origPATH))
	tg, err := Executable(execName)
	s.Require().NoError(err)
	s.Require().Equal(&TelegramDesktop{
		execName,
		s.execPath,
		s.execPath,
	}, tg)
}

func (s *TestTelegramSuite) TestIsSnapFalse() {
	s.CreateFile(s.execPath, true)
	symlinkPath := path.Join(s.dir, symlinkName)
	s.CreateSymlink(s.execPath, symlinkPath)
	tg, err := Executable(symlinkPath)
	s.Require().NoError(err)
	s.Require().False(tg.IsSnap())
}

func (s *TestTelegramSuite) TestIsSnapTrue() {
	snapPath := path.Join(s.dir, "snap")
	s.CreateFile(snapPath, true)
	symlinkPath := path.Join(s.dir, symlinkName)
	s.CreateSymlink(snapPath, symlinkPath)
	tg, err := Executable(symlinkPath)
	s.Require().NoError(err)
	s.Require().True(tg.IsSnap())
}

func TestTelegramSuiteTest(t *testing.T) {
	suite.Run(t, new(TestTelegramSuite))
}
