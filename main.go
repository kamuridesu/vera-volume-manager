package main

import (
	"fmt"

	b "github.com/kamuridesu/vera-volume-manager/internal/bitwarden"
	c "github.com/kamuridesu/vera-volume-manager/internal/config"
	v "github.com/kamuridesu/vera-volume-manager/internal/veracrypt"
)

func main() {
	conf, err := c.LoadConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	bitwarden := b.NewBitwarden(conf.Bitwarden)
	err = bitwarden.Unlock()
	if err != nil {
		panic(err)
	}
	password, err := bitwarden.GetPassword()
	if err != nil {
		panic(err)
	}
	bitwarden.Lock()
	vera := v.NewVeracrypt(conf)
	err = vera.Create(password)
	if err != nil {
		panic(err)
	}
	err = vera.Mount(password)
	if err != nil {
		panic(err)
	}
	c.CreateFolderStructure(conf.DefaultStructure[0], conf.Volume.MountPoint)
	fmt.Println("Press enter to unmount")
	fmt.Scanln()
	vera.Umount()
}
