package utils

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"
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

var ExecuteHook = func(executable string, exitOnFail bool) error {
	fmt.Printf("Executing hook \"%s\"\n", executable)
	cmd := exec.Command("sh", "-c", executable)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		if exitOnFail {
			slog.Error(fmt.Sprintf("hook execution failed: %v", err))
			os.Exit(1)
		}
		return fmt.Errorf("hook execution failed: %w", err)
	}
	return nil
}
