package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	c "github.com/kamuridesu/vera-volume-manager/internal/config"
	"github.com/kamuridesu/vera-volume-manager/internal/keepassxc"
	"github.com/kamuridesu/vera-volume-manager/internal/state"
	v "github.com/kamuridesu/vera-volume-manager/internal/veracrypt"
)

func Check[T any](x T, err error) T {
	CheckErr(err)
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
	fmt.Println("  list     Shows all mounted volumes")
	fmt.Println("\nOptions for all commands (except list):")
	fmt.Println("  -config  Path to the config file (default: ./config.yaml)")
	fmt.Println("  -nohook  Disables hooks after mounting/umounting")
}

func bootstrap(configPath string, ignoreHooks bool) (c.Config, *v.Veracrypt, string) {
	state := Check(state.New())
	conf := Check(c.LoadConfig(configPath, ignoreHooks))
	vera := Check(v.NewVeracrypt(conf, state))

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

	addConfigToCmd := func(cmd *flag.FlagSet) *string {
		return cmd.String("config", "./config.yaml", "Path to config file")
	}

	addNoHookToCmd := func(cmd *flag.FlagSet) *bool {
		return cmd.Bool("nohook", false, "Ignores create/mount/umount hooks")
	}

	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	createConfig := addConfigToCmd(createCmd)
	createNoHook := addNoHookToCmd(createCmd)

	mountCmd := flag.NewFlagSet("mount", flag.ExitOnError)
	mountConfig := addConfigToCmd(mountCmd)
	mountNoHook := addNoHookToCmd(mountCmd)

	umountCmd := flag.NewFlagSet("umount", flag.ExitOnError)
	umountConfig := addConfigToCmd(umountCmd)
	umountAll := umountCmd.Bool("all", false, "Umount all mounted volumes")
	umountNoHook := addNoHookToCmd(umountCmd)

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)

	subcommand := os.Args[1]
	switch subcommand {

	case "create":
		createCmd.Parse(os.Args[2:])
		conf, vera, password := bootstrap(*createConfig, *createNoHook)

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
		_, vera, password := bootstrap(*mountConfig, *mountNoHook)

		fmt.Println("Mounting volume...")
		CheckErr(vera.Mount(password))
		fmt.Println("Volume mounted")

	case "umount":
		umountCmd.Parse(os.Args[2:])

		state := Check(state.New())

		if *umountAll {
			fmt.Println("Umount ALL volumes...")
			v.UmountAll(state)
			fmt.Println("All volumes umounted")
			return
		}

		conf := Check(c.LoadConfig(*umountConfig, *umountNoHook))
		vera := Check(v.NewVeracrypt(conf, state))

		fmt.Println("Unmounting volume...")
		CheckErr(vera.Umount())
		fmt.Println("Volume unmounted")

	case "list":
		listCmd.Parse(os.Args[2:])

		state := Check(state.New())
		mounted := state.GetMountedConfigs()
		if len(mounted) < 1 {
			fmt.Println("No mounted state")
			return
		}

		fmt.Println("Mounted configs: ")
		for _, state := range state.GetMountedConfigs() {
			fmt.Printf("  - %s\n", state)
		}

	default:
		fmt.Printf("Unknown command: %s\n", subcommand)
		printUsage()
		os.Exit(1)
	}
}
