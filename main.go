package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	c "github.com/kamuridesu/vera-volume-manager/internal/config"
	"github.com/kamuridesu/vera-volume-manager/internal/keepassxc"
	v "github.com/kamuridesu/vera-volume-manager/internal/veracrypt"
)

func Check[T any](x T, err error) T {
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	return x
}

func CheckErr(err error) {
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func debug() {
	ss := keepassxc.NewSecretService(c.SecretService{})
	res := Check(ss.GetItem("test", "kamuridesu"))
	fmt.Println(res)
}

func printUsage() {
	scriptName := os.Args[0]
	fmt.Printf("Usage: %s <command> [options]\n", scriptName)
	fmt.Println("\nCommands:")
	fmt.Println("  create   Creates the volume and initializes folder structure")
	fmt.Println("  mount    Mounts the volume")
	fmt.Println("  umount   Unmounts the volume")
	fmt.Println("\nOptions for all commands:")
	fmt.Println("  -config  Path to the config file (default: ./config.yaml)")
}

func bootstrap(configPath string) (c.Config, *v.Veracrypt, string) {
	conf := Check(c.LoadConfig(configPath))
	vera := Check(v.NewVeracrypt(conf))

	ss := keepassxc.NewSecretService(conf.SecretService)
	password := Check(ss.GetPassword())

	return conf, vera, password
}

func main() {
	if os.Getenv("DEBUG") == "1" {
		debug()
		return
	}

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	createConfig := createCmd.String("config", "./config.yaml", "Path to config file")

	mountCmd := flag.NewFlagSet("mount", flag.ExitOnError)
	mountConfig := mountCmd.String("config", "./config.yaml", "Path to config file")

	umountCmd := flag.NewFlagSet("umount", flag.ExitOnError)
	umountConfig := umountCmd.String("config", "./config.yaml", "Path to config file")

	subcommand := os.Args[1]
	switch subcommand {

	case "create":
		createCmd.Parse(os.Args[2:])
		conf, vera, password := bootstrap(*createConfig)

		fmt.Println("Creating volume...")
		CheckErr(vera.Create(password))

		fmt.Println("Mounting to initialize folder structure...")
		CheckErr(vera.Mount(password))

		c.CreateFolderStructure(conf.DefaultStructure, conf.Volume.MountPoint)

		fmt.Println("Unmounting...")
		CheckErr(vera.Umount())
		fmt.Println("Volume created and initialized")

	case "mount":
		mountCmd.Parse(os.Args[2:])
		_, vera, password := bootstrap(*mountConfig)

		fmt.Println("Mounting volume...")
		CheckErr(vera.Mount(password))
		fmt.Println("Volume mounted")

	case "umount":
		umountCmd.Parse(os.Args[2:])

		conf := Check(c.LoadConfig(*umountConfig))
		vera := Check(v.NewVeracrypt(conf))

		fmt.Println("Unmounting volume...")
		CheckErr(vera.Umount())
		fmt.Println("Volume unmounted")

	default:
		fmt.Printf("Unknown command: %s\n", subcommand)
		printUsage()
		os.Exit(1)
	}

}
