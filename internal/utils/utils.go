package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type SeedFile struct {
	Path string
}

func (s *SeedFile) Delete() error {
	return os.Remove(s.Path)
}

func GenerateRandomSeedFile() (*SeedFile, error) {
	if runtime.GOOS == "windows" {
		return &SeedFile{}, nil
	}
	seed := make([]byte, 64)
	_, err := rand.Read(seed)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random values from seed: %w", err)
	}

	seedHex := hex.EncodeToString(seed)
	err = os.WriteFile("SEED.txt", []byte(seedHex), 0644)
	if err != nil {
		return nil, fmt.Errorf("error while saving seed file: %w", err)
	}

	return &SeedFile{Path: "SEED.txt"}, nil
}

type Commands struct {
	create string
	mount  string
	umount string
}

func GetCommands() *Commands {
	if runtime.GOOS == "windows" {
		return &Commands{
			create: "/create %s /password %s /hash sha512 /filesystem %s /size %s /force",
			mount:  "/v %s /l %s /password %s /q",
			umount: "/d %s /q",
		}
	}
	return &Commands{
		create: "-t -c %s --password %s --hash sha512 --filesystem %s --size %s --force --random-source %s --volume-type normal --encryption AES --pim 0 --keyfiles ",
		mount:  "-t --mount %s %s --password %s --pim 0 --protect-hidden no --keyfiles ",
		umount: "-t -d %s",
	}
}

func (c *Commands) Create(volume, password, fs, size string, randomSource string) string {
	if randomSource == "" {
		return fmt.Sprintf(c.create, volume, password, fs, size)
	}
	return fmt.Sprintf(c.create, volume, password, fs, size, randomSource)
}

func (c *Commands) Mount(volume, password, mountPoint string) string {
	return fmt.Sprintf(c.mount, volume, mountPoint, password)
}

func (c *Commands) Umount(volume string) string {
	return fmt.Sprintf(c.umount, volume)
}

func RunCommand(executable string, command string) error {
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

func CreateFolder(folder string) error {
	_, err := os.Stat(folder)
	if err != nil {
		if err := os.MkdirAll(folder, 0755); err != nil {
			return fmt.Errorf("error creating folder: %w", err)
		}
	}
	return nil
}

func ExecuteHook(executable string, exitOnFail bool) error {
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
