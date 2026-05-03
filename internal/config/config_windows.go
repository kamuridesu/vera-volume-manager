//go:build windows

package config

func getValidFileSystems() ([]string, error) {
	return []string{"", "FAT", "ExFAT", "NTFS", "ReFS"}, nil
}
