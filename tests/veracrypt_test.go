package tests

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/kamuridesu/vera-volume-manager/internal/config"
	"github.com/kamuridesu/vera-volume-manager/internal/state"
	"github.com/kamuridesu/vera-volume-manager/internal/utils"
	"github.com/kamuridesu/vera-volume-manager/internal/veracrypt"
	"github.com/stretchr/testify/assert"
)

func setupVeraCryptMock(t *testing.T) (*veracrypt.Veracrypt, *[]string) {
	origRun := utils.RunCommand
	origHook := utils.ExecuteHook
	origFolder := utils.CreateFolder
	origSeed := utils.GenerateRandomSeedFile

	t.Cleanup(func() {
		utils.RunCommand = origRun
		utils.ExecuteHook = origHook
		utils.CreateFolder = origFolder
		utils.GenerateRandomSeedFile = origSeed
	})

	var executedCommands []string

	utils.RunCommand = func(executable string, command string) error {
		executedCommands = append(executedCommands, executable+" "+command)
		return nil
	}
	utils.ExecuteHook = func(executable string, exitOnFail bool) error {
		executedCommands = append(executedCommands, "HOOK: "+executable)
		return nil
	}
	utils.CreateFolder = func(folder string) error { return nil }
	utils.GenerateRandomSeedFile = func() (*utils.SeedFile, error) {
		return &utils.SeedFile{Path: "mock_seed.txt"}, nil
	}

	tmpFile, _ := os.CreateTemp("", "state-*.yaml")
	defer tmpFile.Close()
	mockState := &state.States{States: make(map[string]bool), FilePath: tmpFile.Name()}

	cfg := config.Config{
		File:          "mock_config.yaml",
		VeracryptPath: "/mock/bin",
		Volume: config.Volume{
			Folder:     "/mock/data",
			Name:       "test-vol",
			MountPoint: "/mock/mnt",
			FileSystem: "ExFAT",
			Size:       "10M",
		},
		Hooks: config.Hooks{
			Create: "echo created",
			Mount:  "echo mounted",
			Umount: "echo unmounted",
		},
	}

	vera := &veracrypt.Veracrypt{
		Config:   cfg,
		Commands: utils.GetCommands(),
		States:   mockState,
	}

	return vera, &executedCommands
}

func TestVeracrypt_Mount(t *testing.T) {
	vera, execLog := setupVeraCryptMock(t)

	err := vera.Mount("mysecretpassword")
	assert.NoError(t, err)

	assert.Len(t, *execLog, 2, "Should have run Veracrypt and the hook")

	expectedVolPath := filepath.Join("/mock/data", "test-vol")
	assert.Contains(t, (*execLog)[0], "/mock/bin/veracrypt")
	assert.Contains(t, (*execLog)[0], expectedVolPath)
	assert.Contains(t, (*execLog)[0], "--password mysecretpassword")
	assert.Contains(t, (*execLog)[0], "/mock/mnt")

	assert.Equal(t, "HOOK: echo mounted", (*execLog)[1])

	absPath, err := filepath.Abs(vera.Config.File)
	assert.NoError(t, err)
	assert.True(t, vera.States.States[absPath], "State should be saved using the absolute path")
}

func TestVeracrypt_Create(t *testing.T) {
	vera, execLog := setupVeraCryptMock(t)

	err := vera.Create("newpassword")
	assert.NoError(t, err)

	assert.Len(t, *execLog, 2, "Should have run Veracrypt and the hook")

	expectedVolPath := filepath.Join("/mock/data", "test-vol")
	assert.Contains(t, (*execLog)[0], "/mock/bin/veracrypt")
	assert.Contains(t, (*execLog)[0], expectedVolPath)
	assert.Contains(t, (*execLog)[0], "--password newpassword")
	assert.Contains(t, (*execLog)[0], "ExFAT")
	assert.Contains(t, (*execLog)[0], "10M")
	assert.Contains(t, (*execLog)[0], "mock_seed.txt")

	assert.Equal(t, "HOOK: echo created", (*execLog)[1])

	absPath, err := filepath.Abs(vera.Config.File)
	assert.NoError(t, err)

	stateVal, exists := vera.States.States[absPath]
	assert.True(t, exists, "State entry should exist")
	assert.False(t, stateVal, "State should be saved as unmounted (false)")
}

func TestVeracrypt_Umount(t *testing.T) {
	vera, execLog := setupVeraCryptMock(t)

	err := vera.Umount()
	assert.NoError(t, err)

	assert.Len(t, *execLog, 2, "Should have run the hook and Veracrypt")

	assert.Equal(t, "HOOK: echo unmounted", (*execLog)[0])

	assert.Contains(t, (*execLog)[1], "/mock/bin/veracrypt")
	assert.Contains(t, (*execLog)[1], "/mock/mnt")

	absPath, err := filepath.Abs(vera.Config.File)
	assert.NoError(t, err)

	stateVal, exists := vera.States.States[absPath]
	assert.True(t, exists, "State entry should exist")
	assert.False(t, stateVal, "State should be saved as unmounted (false)")
}

func TestNewVeracrypt_Success(t *testing.T) {
	tmpDir := t.TempDir()
	exeName := "veracrypt"
	if runtime.GOOS == "windows" {
		exeName = "VeraCrypt.exe"
	}

	os.WriteFile(filepath.Join(tmpDir, exeName), []byte("dummy binary"), 0755)

	cfg := config.Config{VeracryptPath: tmpDir}
	vc, err := veracrypt.NewVeracrypt(cfg, &state.States{})

	assert.NoError(t, err)
	assert.NotNil(t, vc)
}

func TestNewVeracrypt_ValidationFail(t *testing.T) {
	cfg := config.Config{VeracryptPath: "/fake/path/that/does/not/exist"}
	_, err := veracrypt.NewVeracrypt(cfg, &state.States{})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get path")
}
