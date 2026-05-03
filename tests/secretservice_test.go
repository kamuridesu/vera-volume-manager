package tests

import (
	"errors"
	"testing"

	"github.com/kamuridesu/vera-volume-manager/internal/config"
	"github.com/kamuridesu/vera-volume-manager/internal/keepassxc"
	"github.com/stretchr/testify/assert"
)

func TestSecretService_GetPassword_Success(t *testing.T) {
	keepassxc.KeyringGet = func(service, user string) (string, error) {
		assert.Equal(t, "test-service", service)
		assert.Equal(t, "kamuridesu", user)
		return "superSecretPassword123", nil
	}

	cfg := config.SecretService{
		ServiceName: "test-service",
		Username:    "kamuridesu",
	}
	ss := keepassxc.NewSecretService(cfg)

	pass, err := ss.GetPassword()
	assert.NoError(t, err)
	assert.Equal(t, "superSecretPassword123", pass)
}

func TestSecretService_GetPassword_Error(t *testing.T) {
	keepassxc.KeyringGet = func(service, user string) (string, error) {
		return "", errors.New("secret not found")
	}

	cfg := config.SecretService{
		ServiceName: "bad-service",
		Username:    "nobody",
	}
	ss := keepassxc.NewSecretService(cfg)

	pass, err := ss.GetPassword()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error retrieving password")
	assert.Empty(t, pass)
}
