package veracrypt

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	c "github.com/kamuridesu/vera-volume-manager/internal/config"
	u "github.com/kamuridesu/vera-volume-manager/internal/utils"
)

type Veracrypt struct {
	Config   c.Config
	Commands *u.Commands
}

func NewVeracrypt(config c.Config) *Veracrypt {
	return &Veracrypt{
		Config:   config,
		Commands: u.GetCommands(),
	}
}

func (v *Veracrypt) getCreateExecutablePath() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(v.Config.VeracryptPath, "VeraCrypt Format.exe")
	}
	return filepath.Join(v.Config.VeracryptPath, "veracrypt")
}

func (v *Veracrypt) getManageExecutablePath() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(v.Config.VeracryptPath, "VeraCrypt.exe")
	}
	return filepath.Join(v.Config.VeracryptPath, "veracrypt")
}

func (v *Veracrypt) Create(password string) error {
	executable := v.getCreateExecutablePath()
	randomFile, err := u.GenerateRandomSeedFile()
	if err != nil {
		return err
	}
	defer randomFile.Delete()
	_, err = os.Stat(v.Config.Volume.Folder)
	if err != nil {
		if err := os.MkdirAll(v.Config.Volume.Folder, 0755); err != nil {
			return err
		}
	}
	targetPath := filepath.Join(v.Config.Volume.Folder, v.Config.Volume.Name)
	command := v.Commands.Create(targetPath, password, v.Config.Volume.Size, randomFile.Path)
	if err := u.RunCommand(executable, command); err != nil {
		return err
	}
	fmt.Println("Volume created at", targetPath)
	return nil
}

func (v *Veracrypt) Mount(password string) error {
	executable := v.getManageExecutablePath()
	command := v.Commands.Mount(filepath.Join(v.Config.Volume.Folder, v.Config.Volume.Name), password, v.Config.Volume.MountPoint)
	if err := u.RunCommand(executable, command); err != nil {
		return err
	}
	fmt.Println("Volume mounted at", v.Config.Volume.MountPoint)
	return nil
}

func (v *Veracrypt) Umount() error {
	executable := v.getManageExecutablePath()
	command := v.Commands.Umount(v.Config.Volume.MountPoint)
	if err := u.RunCommand(executable, command); err != nil {
		return err
	}
	fmt.Println("Volume unmounted")
	return nil
}
