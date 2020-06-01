package xdg

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type BaseSuite struct {
	suite.Suite
	varName      string
	defaultValue string
	origVarIsSet bool
	origVarValue string
}

func (s *BaseSuite) SetupTest() {
	s.Require().NotEmpty(s.varName)
	s.Require().NotEmpty(s.defaultValue)
	value, found := os.LookupEnv(s.varName)
	s.origVarValue = value
	s.origVarIsSet = found
}

func (s *BaseSuite) TearDownTest() {
	if s.origVarIsSet {
		os.Setenv(s.varName, s.origVarValue)
	} else {
		os.Unsetenv(s.varName)
	}
}

func (s *BaseSuite) Set(value string) {
	os.Setenv(s.varName, value)
}

func (s *BaseSuite) Unset() {
	os.Unsetenv(s.varName)
}

func (s *BaseSuite) TestNoVar() {
	s.Unset()
	s.Require().Equal(s.defaultValue, GetConfigHome())
}

func (s *BaseSuite) TestEmptyVar() {
	s.Set("")
	s.Require().Equal(s.defaultValue, GetConfigHome())
}

func (s *BaseSuite) TestRelPath() {
	s.Set("not/allowed")
	s.Require().Equal(s.defaultValue, GetConfigHome())
}

func (s *BaseSuite) TestAbsPath() {
	s.Set("/some/path")
	s.Require().Equal("/some/path", GetConfigHome())
}

type TestGetConfigHomeSuite struct {
	BaseSuite
}

func (s *TestGetConfigHomeSuite) SetupSuite() {
	s.varName = "XDG_CONFIG_HOME"
	s.defaultValue = os.ExpandEnv("$HOME/.config")
}

func TestGetConfigHomeSuiteTest(t *testing.T) {
	suite.Run(t, new(TestGetConfigHomeSuite))
}
