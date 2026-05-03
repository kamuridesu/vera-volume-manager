//go:build windows

package utils

func GetCommands() *Commands {
	return &Commands{
		create: "/create %s /password %s /hash sha512 /filesystem %s /size %s /force",
		mount:  "/v %s /l %s /password %s /q",
		umount: "/d %s /q",
	}
}

func ElevateCommand(executable string, command string) (string, string) {
	return executable, command
}
