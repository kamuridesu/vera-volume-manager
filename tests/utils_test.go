package tests

import (
	"os"
	"testing"

	"github.com/kamuridesu/vera-volume-manager/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestCommands_Create(t *testing.T) {
	cmds := utils.GetCommands()

	cmdStr := cmds.Create("/path/to/vol", "mypass", "ExFAT", "10M", "SEED.txt")

	assert.Contains(t, cmdStr, "/path/to/vol")
	assert.Contains(t, cmdStr, "mypass")
	assert.Contains(t, cmdStr, "ExFAT")
	assert.Contains(t, cmdStr, "10M")
}

func TestGenerateRandomSeedFile(t *testing.T) {
	seed, err := utils.GenerateRandomSeedFile()
	assert.NoError(t, err)

	if seed.Path != "" {
		_, err = os.Stat(seed.Path)
		assert.NoError(t, err, "Seed file should be created on disk")

		err = seed.Delete()
		assert.NoError(t, err, "Should delete seed file without error")

		_, err = os.Stat(seed.Path)
		assert.Error(t, err, "Seed file should no longer exist")
	}
}
