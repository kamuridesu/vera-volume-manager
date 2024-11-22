package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Volume struct {
	Folder     string `yaml:"folder"`
	Name       string `yaml:"name"`
	MountPoint string `yaml:"mount_point"`
	Size       string `yaml:"size"`
}

type Bitwarden struct {
	Url            string `yaml:"url"`
	PasswordB64    string `yaml:"password_base64"`
	CredentialName string `yaml:"credential_name"`
}

type Folder struct {
	Name     string   `yaml:"name"`
	Children []Folder `yaml:"children"`
}

type Config struct {
	VeracryptPath    string    `yaml:"veracrypt_path"`
	Volume           Volume    `yaml:"volume"`
	DefaultStructure []Folder  `yaml:"default_structure"`
	Bitwarden        Bitwarden `yaml:"bitwarden"`
}

func LoadConfig() (Config, error) {
	file, err := os.Open("config.yaml")
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

	return config, nil
}

func CreateFolderStructure(folder Folder, parent string) {
	if parent == "" {
		parent = "."
	}
	folderName := filepath.Join(parent, folder.Name)
	fmt.Println("Creating folder", folderName)
	err := os.Mkdir(folderName, 0755)
	if err != nil {
		fmt.Println(err)
	}
	for _, child := range folder.Children {
		CreateFolderStructure(child, folderName)
	}
}
