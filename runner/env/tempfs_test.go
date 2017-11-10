package env

import (
	"testing"

	"os"
	"path/filepath"

	"github.com/bivas/rivi/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetWithoutCreate(t *testing.T) {
	tested := tempFSEnvironmentProvider()
	assert.Empty(t, tested.ConfigFilePath(), "no config file")
}

func TestGetWithCreate(t *testing.T) {
	tested := tempFSEnvironmentProvider()
	data := &mocks.MockData{
		RulesFileContent: "content",
	}
	assert.Nil(t, tested.Create(data), "create should not fail")
	assert.NotEmpty(t, tested.ConfigFilePath(), "no config file")
}

func TestCleanup(t *testing.T) {
	tested := tempFSEnvironmentProvider()
	data := &mocks.MockData{
		RulesFileContent: "content",
	}
	assert.Nil(t, tested.Create(data), "create should not fail")
	tempdir := filepath.Dir(tested.ConfigFilePath())
	_, err := os.Stat(tempdir)
	assert.Nil(t, err, "should exist")
	tested.Cleanup()
	_, err = os.Stat(tempdir)
	assert.True(t, os.IsNotExist(err), "should be cleaned")
}

func TestDefaultProvider(t *testing.T) {
	tested, err := GetEnvironment()
	assert.Nil(t, err, "should create env")
	_, ok := tested.(*tempFSEnvironment)
	assert.True(t, ok, "default is temp fs")
}
