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

func NewVeracrypt(config c.Config) (*Veracrypt, error) {
	vera := &Veracrypt{
		Config:   config,
		Commands: u.GetCommands(),
	}
	return vera.validate()
}

func (v *Veracrypt) validate() (*Veracrypt, error) {
	for _, path := range []string{v.getManageExecutablePath(), v.getManageExecutablePath()} {
		if _, err := os.Stat(path); err != nil {
			return nil, fmt.Errorf("failed to get path for veracrypt executable: %w", err)
		}
	}
	return v, nil
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
		return fmt.Errorf("failed to generate random seed file: %w", err)
	}
	defer randomFile.Delete()
	if err := u.CreateFolder(v.Config.Volume.Folder); err != nil {
		return err
	}
	if err := u.CreateFolder(v.Config.Volume.MountPoint); err != nil {
		return fmt.Errorf("err creating mount point folder: %w", err)
	}

	targetPath := filepath.Join(v.Config.Volume.Folder, v.Config.Volume.Name)
	command := v.Commands.Create(targetPath, password, v.Config.Volume.FileSystem, v.Config.Volume.Size, randomFile.Path)
	// fmt.Println(executable + " " + command)
	if err := u.RunCommand(executable, command); err != nil {
		return fmt.Errorf("error running command: %w", err)
	}
	// fmt.Println("Volume created at", targetPath)
	u.ExecuteHook(v.Config.Hooks.Create, v.Config.Hooks.ExitOnFailed)
	return nil
}

func (v *Veracrypt) Mount(password string) error {
	executable := v.getManageExecutablePath()
	command := v.Commands.Mount(filepath.Join(v.Config.Volume.Folder, v.Config.Volume.Name), password, v.Config.Volume.MountPoint)
	if err := u.RunCommand(executable, command); err != nil {
		// fmt.Println(executable + " " + command)
		return err
	}
	fmt.Println("Volume mounted at", v.Config.Volume.MountPoint)
	u.ExecuteHook(v.Config.Hooks.Mount, v.Config.Hooks.ExitOnFailed)
	return nil
}

func (v *Veracrypt) Umount() error {
	executable := v.getManageExecutablePath()
	command := v.Commands.Umount(v.Config.Volume.MountPoint)
	if err := u.RunCommand(executable, command); err != nil {
		return err
	}
	fmt.Println("Volume unmounted")
	u.ExecuteHook(v.Config.Hooks.Umount, v.Config.Hooks.ExitOnFailed)
	return nil
}
