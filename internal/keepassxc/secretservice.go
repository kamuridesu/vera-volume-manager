package keepassxc

import (
	"fmt"

	c "github.com/kamuridesu/vera-volume-manager/internal/config"
	"github.com/zalando/go-keyring"
)

type SecretService struct {
	Config c.SecretService
}

func NewSecretService(config c.SecretService) *SecretService {
	return &SecretService{Config: config}
}

func (s *SecretService) GetItem(service string, user string) (string, error) {
	password, err := keyring.Get(service, user)

	if err != nil {
		return "", fmt.Errorf("error retrieving password: %w", err)
	}

	return password, err
}

func (s *SecretService) GetPassword() (string, error) {
	return s.GetItem(s.Config.ServiceName, s.Config.Username)
}
