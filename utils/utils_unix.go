//go:build !windows

package utils

import (
	"fmt"
	"os"
	"syscall"
)

// GetHardLinkID returns the string representing unique hardlink identifier
func GetHardLinkID(path string) (string, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	stat, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		return "", fmt.Errorf("not a unix stat struct")
	}
	return fmt.Sprintf("%d_%d", stat.Dev, stat.Ino), nil
}
