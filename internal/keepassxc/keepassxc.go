package keepassxc

import (
	"os"
	"os/exec"

	c "github.com/kamuridesu/vera-volume-manager/internal/config"
	u "github.com/kamuridesu/vera-volume-manager/internal/utils"
)

type Keepassxc struct {
	Config *c.Keepassxc
}

func NewKeepassxc(config *c.Keepassxc) (*Keepassxc, error) {
	if _, err := os.Stat(config.File); err != nil {
		return nil, err
	}
	return &Keepassxc{Config: config}, nil
}

func (k *Keepassxc) GetPassword() (string, error) {
	cmd := exec.Command("keepassxc-cli", "show", k.Config.File, k.Config.CredentialPath)
	cmd.Stdin
}
