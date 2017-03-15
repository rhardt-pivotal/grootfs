package filesystems

import (
	"syscall"

	errorspkg "github.com/pkg/errors"
)

func CheckFSPath(path string, expectedFilesystem int64, expectedFilesystemName string) error {
	statfs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &statfs)
	if err != nil {
		return errorspkg.Wrapf(err, "Failed to detect type of filesystem")
	}

	if statfs.Type != expectedFilesystem {
		return errorspkg.Errorf("filesystem driver requires store filesystem to be %s", expectedFilesystemName)
	}
	return nil
}