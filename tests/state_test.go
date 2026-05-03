package tests

import (
	"os"
	"testing"

	"github.com/kamuridesu/vera-volume-manager/internal/state"
	"github.com/stretchr/testify/assert"
)

func TestSaveState(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "state-*.yaml")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	s := &state.States{
		States:   make(map[string]bool),
		FilePath: tmpFile.Name(),
	}

	err = s.SaveState("/my/config/path.yaml", true)
	assert.NoError(t, err)

	val, exists := s.States["/my/config/path.yaml"]
	assert.True(t, exists)
	assert.True(t, val)

	content, err := os.ReadFile(tmpFile.Name())
	assert.NoError(t, err)
	assert.Contains(t, string(content), "/my/config/path.yaml: true")
}
