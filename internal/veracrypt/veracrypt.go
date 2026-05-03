package veracrypt

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"

	c "github.com/kamuridesu/vera-volume-manager/internal/config"
	"github.com/kamuridesu/vera-volume-manager/internal/state"
	u "github.com/kamuridesu/vera-volume-manager/internal/utils"
)

type Veracrypt struct {
	Config   c.Config
	Commands *u.Commands
	States   *state.States
}

func NewVeracrypt(config c.Config, states *state.States) (*Veracrypt, error) {
	vera := &Veracrypt{
		Config:   config,
		Commands: u.GetCommands(),
		States:   states,
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
	if err := u.CreateFolder(v.Config.Volume.Folder); err != nil {
		return err
	}
	if err := u.CreateFolder(v.Config.Volume.MountPoint); err != nil {
		return fmt.Errorf("err creating mount point folder: %w", err)
	}

	targetPath := filepath.Join(v.Config.Volume.Folder, v.Config.Volume.Name)
	command := v.Commands.Create(targetPath, password, v.Config.Volume.FileSystem, v.Config.Volume.Size)
	execCmd, execArgs := u.ElevateCommand(executable, command)
	// fmt.Println(executable + " " + command)
	if err := u.RunCommand(execCmd, execArgs); err != nil {
		return fmt.Errorf("error running command: %w", err)
	}
	// fmt.Println("Volume created at", targetPath)
	v.States.SaveState(v.Config.File, false)
	u.ExecuteHook(v.Config.Hooks.Create, v.Config.Hooks.ExitOnFailed)
	return nil
}

func (v *Veracrypt) Mount(password string) error {
	executable := v.getManageExecutablePath()
	command := v.Commands.Mount(filepath.Join(v.Config.Volume.Folder, v.Config.Volume.Name), password, v.Config.Volume.MountPoint)
	execCmd, execArgs := u.ElevateCommand(executable, command)
	if err := u.RunCommand(execCmd, execArgs); err != nil {
		return fmt.Errorf("error running command: %w", err)
	}
	v.States.SaveState(v.Config.File, true)

	u.ExecuteHook(v.Config.Hooks.Mount, v.Config.Hooks.ExitOnFailed)
	return nil
}

func (v *Veracrypt) Umount() error {
	u.ExecuteHook(v.Config.Hooks.Umount, v.Config.Hooks.ExitOnFailed)
	executable := v.getManageExecutablePath()
	command := v.Commands.Umount(v.Config.Volume.MountPoint)
	execCmd, execArgs := u.ElevateCommand(executable, command)
	if err := u.RunCommand(execCmd, execArgs); err != nil {
		return fmt.Errorf("error running command: %w", err)
	}
	v.States.SaveState(v.Config.File, false)
	return nil
}

func UmountAll(s *state.States) {
	for config, isMounted := range s.States {
		if !isMounted {
			continue
		}
		cfg, err := c.LoadConfig(config)
		if err != nil {
			slog.Error(fmt.Sprintf("failed to get config for %s: %s", config, err))
			continue
		}
		vera, err := NewVeracrypt(cfg, s)
		if err != nil {
			slog.Error(fmt.Sprintf("veracrypt failed: %s", err))
			continue
		}
		err = vera.Umount()
		if err != nil {
			slog.Error("failed to umount")
			continue
		}
	}
}
