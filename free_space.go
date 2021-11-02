package util

import "syscall"

// FreeSpace uses a syscall to get the amount of free bytes in a directory
func FreeSpace(dir string) (uint64, error) {
	var statfs syscall.Statfs_t
	err := syscall.Statfs(dir, &statfs)
	if err != nil {
		return 0, err
	}

	return statfs.Bavail * uint64(statfs.Bsize), nil
}

// DiskSpace uses a syscall to get the size of a storage device
func DiskSpace(dir string) (uint64, error) {
	var statfs syscall.Statfs_t
	err := syscall.Statfs(dir, &statfs)
	if err != nil {
		return 0, err
	}

	return statfs.Blocks * uint64(statfs.Bsize), nil
}
