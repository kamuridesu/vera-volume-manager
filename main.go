package main

import (
	"flag"
	"fmt"
	"os"

	b "github.com/kamuridesu/vera-volume-manager/internal/bitwarden"
	c "github.com/kamuridesu/vera-volume-manager/internal/config"
	v "github.com/kamuridesu/vera-volume-manager/internal/veracrypt"
)

var (
	ConfigFileLocation = "./config.yaml"
)

func argparse() {
	flag.Parse()
	args := flag.Args()
	if len(args) > 0 {
		ConfigFileLocation = args[0]
	}
}

func exitOnError(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func main() {
	argparse()
	conf, err := c.LoadConfig(ConfigFileLocation)
	if err != nil {
		exitOnError(err)
	}
	bitwarden, err := b.NewBitwarden(conf.Bitwarden)
	if err != nil {
		fmt.Println("Fail to create Bitwarden client, check if client is running. Error is:", err)
		return
	}
	err = bitwarden.Unlock()
	if err != nil {
		exitOnError(err)
	}
	password, err := bitwarden.GetPassword()
	if err != nil {
		exitOnError(err)
	}
	bitwarden.Lock()
	vera, err := v.NewVeracrypt(conf)
	if err != nil {
		fmt.Printf("Could not find any Veracrypt executable in %s, error is: %v\n", conf.VeracryptPath, err)
		return
	}
	err = vera.Create(password)
	if err != nil {
		exitOnError(err)
	}
	err = vera.Mount(password)
	if err != nil {
		exitOnError(err)
	}
	c.CreateFolderStructure(conf.DefaultStructure, conf.Volume.MountPoint)
	fmt.Println("Press enter to unmount")
	fmt.Scanln()
	vera.Umount()
}
