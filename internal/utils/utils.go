package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func DecodeBase64String(base64String string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

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
		return nil, err
	}

	seedHex := hex.EncodeToString(seed)
	err = os.WriteFile("SEED.txt", []byte(seedHex), 0644)
	if err != nil {
		return nil, err
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
		mount:  "-t --mount %s %s --password %s --pim 0 --protect-hidden no --slot 1 --keyfiles ",
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
	if err := cmd.Start(); err != nil {
		return err
	}
	cmd.Wait()
	if cmd.ProcessState.ExitCode() != 0 {
		return fmt.Errorf("command failed with exit code %d", cmd.ProcessState.ExitCode())
	}
	return nil
}
