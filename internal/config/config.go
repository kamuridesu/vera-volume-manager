package config

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/goccy/go-yaml"
)

type HookType string

const (
	Create HookType = "create"
	Mount  HookType = "mount"
	Umount HookType = "umount"
)

type Volume struct {
	Folder     string `yaml:"folder"`
	Name       string `yaml:"name"`
	MountPoint string `yaml:"mount_point"`
	Size       string `yaml:"size"`
	FileSystem string `yaml:"filesystem"`
}

type Hooks struct {
	Create       string `yaml:"create"`
	Mount        string `yaml:"mount"`
	Umount       string `yaml:"umount"`
	ExitOnFailed bool   `yaml:"exit_on_failed"`
}

type SecretService struct {
	ServiceName string `yaml:"service_name"`
	Username    string `yaml:"username"`
}

type Folder struct {
	Name     string   `yaml:"name"`
	Children []Folder `yaml:"children"`
}

type Config struct {
	VeracryptPath    string        `yaml:"veracrypt_path"`
	Volume           Volume        `yaml:"volume"`
	DefaultStructure []Folder      `yaml:"default_structure"`
	SecretService    SecretService `yaml:"secret_service"`
	Hooks            Hooks         `yaml:"hooks"`
	File             string
	IgnoreHooks      bool
}

func LoadConfig(filename string, ignoreHooks bool) (Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	config := Config{}
	err = decoder.Decode(&config)
	if err != nil {
		return Config{}, err
	}

	valid_fs, err := getValidFileSystems()
	if err != nil {
		return Config{}, err
	}

	if config.Volume.FileSystem == "" {
		config.Volume.FileSystem = "ExFAT"
	}

	if !slices.Contains(valid_fs, config.Volume.FileSystem) {
		return Config{}, fmt.Errorf("variable 'filesystem' can only be %s, got '%s'", strings.Join(valid_fs, ", "), config.Volume.FileSystem)
	}

	config.File = filename
	config.IgnoreHooks = ignoreHooks

	replaceHookVariables(&config)

	return config, nil
}

func replaceHookVariables(config *Config) {
	replacements := map[string]string{
		".volume.folder":      filepath.Clean(config.Volume.Folder),
		".volume.name":        config.Volume.Name,
		".volume.mount_point": filepath.Clean(config.Volume.MountPoint),
		".volume.size":        config.Volume.Size,
		".volume.filesystem":  config.Volume.FileSystem,
		".veracrypt_path":     filepath.Clean(config.VeracryptPath),
	}

	replace := func(cmd string) string {
		resolved := cmd
		for placeholder, value := range replacements {
			resolved = strings.ReplaceAll(resolved, placeholder, value)
		}
		return resolved
	}

	config.Hooks.Create = replace(config.Hooks.Create)
	config.Hooks.Mount = replace(config.Hooks.Mount)
	config.Hooks.Umount = replace(config.Hooks.Umount)
}

func CreateFolderStructure(folder []Folder, parent string) {
	if parent == "" {
		parent = "."
	}
	for _, child := range folder {
		folderName := filepath.Join(parent, child.Name)
		fmt.Println("Creating folder", folderName)
		err := os.MkdirAll(folderName, 0755)
		if err != nil {
			fmt.Println(err)
		}
		CreateFolderStructure(child.Children, folderName)
	}
}
