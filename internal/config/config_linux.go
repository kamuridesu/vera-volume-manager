//go:build linux

package config

func getValidFileSystems() ([]string, error) {
	return []string{"", "FAT", "ExFAT", "NTFS", "Ext2", "Ext3", "Ext4", "Btrfs"}, nil
}
