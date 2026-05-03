package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kamuridesu/vera-volume-manager/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_Valid(t *testing.T) {
	yamlContent := `
veracrypt_path: /usr/bin/
volume:
  folder: ./test-folder
  name: test-vol
  mount_point: /mnt/test
  size: 50M
  filesystem: FAT
hooks:
  create: echo "creating .volume.name"
`
	tmpFile, err := os.CreateTemp("", "config-*.yaml")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write([]byte(yamlContent))
	assert.NoError(t, err)
	tmpFile.Close()

	cfg, err := config.LoadConfig(tmpFile.Name())
	assert.NoError(t, err)

	assert.Equal(t, "/usr/bin/", cfg.VeracryptPath)
	assert.Equal(t, "FAT", cfg.Volume.FileSystem)
	assert.Equal(t, "echo \"creating test-vol\"", cfg.Hooks.Create, "Hook variable should be replaced")
}

func TestLoadConfig_DefaultFileSystem(t *testing.T) {
	yamlContent := `
volume:
  folder: ./test-folder
  name: test-vol
`
	tmpFile, err := os.CreateTemp("", "config-*.yaml")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	os.WriteFile(tmpFile.Name(), []byte(yamlContent), 0644)

	cfg, err := config.LoadConfig(tmpFile.Name())
	assert.NoError(t, err)

	assert.Equal(t, "ExFAT", cfg.Volume.FileSystem, "Empty filesystem should default to ExFAT")
}

func TestLoadConfig_InvalidFileSystem(t *testing.T) {
	yamlContent := `
volume:
  filesystem: ext4
`
	tmpFile, err := os.CreateTemp("", "config-*.yaml")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	os.WriteFile(tmpFile.Name(), []byte(yamlContent), 0644)

	_, err = config.LoadConfig(tmpFile.Name())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "variable 'filesystem' can be only FAT or ExFAT")
}

func TestCreateFolderStructure(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vvm-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	structure := []config.Folder{
		{
			Name: "docs",
			Children: []config.Folder{
				{Name: "pdf"},
			},
		},
	}

	config.CreateFolderStructure(structure, tmpDir)

	_, err = os.Stat(filepath.Join(tmpDir, "docs"))
	assert.NoError(t, err, "docs directory should exist")

	_, err = os.Stat(filepath.Join(tmpDir, "docs", "pdf"))
	assert.NoError(t, err, "docs/pdf directory should exist")
}
