//go:build !windows

package utils

import "os"

func GetCommands() *Commands {
	return &Commands{
		create: "-t -c %s --password %s --hash sha512 --filesystem %s --size %s --force --random-source /dev/urandom --volume-type normal --encryption AES --pim 0 --keyfiles ",
		mount:  "-t --mount %s %s --password %s --pim 0 --protect-hidden no --keyfiles ",
		umount: "-t -d %s",
	}
}

func ElevateCommand(executable string, command string) (string, string) {
	if os.Geteuid() != 0 {
		return "sudo", executable + " " + command
	}
	return executable, command
}
