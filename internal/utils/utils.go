package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/kamuridesu/vera-volume-manager/internal/config"
)

type SeedFile struct {
	Path string
}

func (s *SeedFile) Delete() error {
	return os.Remove(s.Path)
}

type Commands struct {
	create string
	mount  string
	umount string
}

func (c *Commands) Create(volume, password, fs, size string) string {
	return fmt.Sprintf(c.create, volume, password, fs, size)
}

func (c *Commands) Mount(volume, password, mountPoint string) string {
	return fmt.Sprintf(c.mount, volume, mountPoint, password)
}

func (c *Commands) Umount(volume string) string {
	return fmt.Sprintf(c.umount, volume)
}

var RunCommand = func(executable string, command string) error {
	cmd := exec.Command(executable, strings.Split(command, " ")...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error while executing command %s with args %s: %w", executable, command, err)
	}
	cmd.Wait()
	if cmd.ProcessState.ExitCode() != 0 {
		return fmt.Errorf("command failed with exit code %d", cmd.ProcessState.ExitCode())
	}
	return nil
}

var CreateFolder = func(folder string) error {
	_, err := os.Stat(folder)
	if err != nil {
		if err := os.MkdirAll(folder, 0755); err != nil {
			return fmt.Errorf("error creating folder: %w", err)
		}
	}
	return nil
}

var ExecuteHook = func(cfg *config.Config, hookType config.HookType) error {
	if cfg.IgnoreHooks == true {
		fmt.Printf("Ignoring hook %s\n", string(hookType))
		return nil
	}

	switch hookType {
	case config.Create:
		return runHook(string(cfg.Hooks.Create), cfg.Hooks.ExitOnFailed)
	case config.Mount:
		return runHook(string(cfg.Hooks.Mount), cfg.Hooks.ExitOnFailed)
	case config.Umount:
		return runHook(string(cfg.Hooks.Umount), cfg.Hooks.ExitOnFailed)
	default:
		return nil
	}
}
