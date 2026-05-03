package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kamuridesu/vera-volume-manager/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestCommands_Create(t *testing.T) {
	cmds := utils.GetCommands()

	cmdStr := cmds.Create("/path/to/vol", "mypass", "ExFAT", "10M")

	assert.Contains(t, cmdStr, "/path/to/vol")
	assert.Contains(t, cmdStr, "mypass")
	assert.Contains(t, cmdStr, "ExFAT")
	assert.Contains(t, cmdStr, "10M")
}

func TestRealUtils_CreateFolder(t *testing.T) {
	tmpDir := t.TempDir()
	target := filepath.Join(tmpDir, "new-sub-folder")

	err := utils.CreateFolder(target)
	assert.NoError(t, err)

	_, err = os.Stat(target)
	assert.NoError(t, err)

	err = utils.CreateFolder(target)
	assert.NoError(t, err)
}

func TestRealUtils_RunCommand(t *testing.T) {
	err := utils.RunCommand("go", "version")
	assert.NoError(t, err)

	err = utils.RunCommand("some-command-that-definitely-doesnt-exist", "")
	assert.Error(t, err)
}

func TestRealUtils_ExecuteHook(t *testing.T) {
	err := utils.ExecuteHook("echo 'hook works'", false)
	assert.NoError(t, err)

	err = utils.ExecuteHook("false", false)
	assert.Error(t, err, "Should fail when command fails")
}
