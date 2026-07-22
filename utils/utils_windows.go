//go:build windows

package utils

import (
	"fmt"
	"syscall"
)

// GetHardLinkID returns the string representing unique hardlink identifier
func GetHardLinkID(path string) (string, error) {
	fd, err := syscall.Open(path, syscall.O_RDONLY, 0)
	if err != nil {
		return "", err
	}
	defer syscall.Close(fd)

	var info syscall.ByHandleFileInformation
	err = syscall.GetFileInformationByHandle(fd, &info)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d_%d_%d",
		info.VolumeSerialNumber,
		info.FileIndexHigh,
		info.FileIndexLow,
	), nil
}
