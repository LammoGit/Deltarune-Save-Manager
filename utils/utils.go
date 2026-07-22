// Package utils holds utility functions used by other packages in the project
package utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"slices"
)

// MAXCHAPTER stores maximum available chapter
const MAXCHAPTER = 5

// variables storing errors
var (
	ErrOSNotSupported         = errors.New("your OS is not supported")
	ErrOSNotSupportedHardLink = fmt.Errorf("%w due to hard links handling", ErrOSNotSupported)
	ErrEmptySaveName          = errors.New("empty save name is passed which is not allowed")
	ErrSaveNameIsTaken        = errors.New("given save name is already taken")
	ErrChapterNotSupported    = errors.New("given chapter is not yet supported")
	ErrSaveNotExist           = errors.New("at least on of the given saves doesn't exist")
	ErrShortSaveFile          = errors.New("save file doesn't contain enough lines for parsing")
	ErrValueCannotBeSet       = errors.New("can't set value during save parsing (Internal error)")
	ErrWrongLineType          = errors.New("save file field contains value of an invalid type")
	ErrTakenByUnmanagedSave   = errors.New("slot is taken by an unmanaged save")
)

// variables storing compiled regex patterns
var (
	SlotSectionLabelRegex = regexp.MustCompile(`G(?:[2-7]_)?\d+`)
	SlotRegex             = regexp.MustCompile(`filech(0|[^0\D]\d*)_(0|[^0\D]\d*)(_b)?`)
	SaveRegex             = regexp.MustCompile(`(0|[^0\D]\d*)_(a|b)_(.+)`)
)

// GetSavesFolderPath returns Deltarune save folder path depending on user's OS
func GetSavesFolderPath() (string, error) {
	switch runtime.GOOS {
	case "windows":
		dirPath, err := os.UserCacheDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(dirPath, "DELTARUNE"), nil
	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, "Library", "Application Support", "com.tobyfox.deltarune"), nil
	case "linux":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, ".steam", "steam", "steamapps", "compatdata", "1671210", "pfx", "drive_c", "users", "steamuser", "AppData", "Local", "DELTARUNE"), nil
	default:
		return "", fmt.Errorf("OS not supported: %s", runtime.GOOS)
	}
}

// DeleteEqual deletes all instances of an element from a slice
func DeleteEqual[S ~[]E, E comparable](s S, el E) S {
	return slices.DeleteFunc(s, func(x E) bool { return x == el })
}

// FileExists checks existence of a file
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// DirExists checks existence of a directory
func DirExists(dirPath string) bool {
	info, err := os.Stat(dirPath)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// TempFilePath returns a random path for a temporary in the given directory
func TempFilePath(dirPath string) (string, error) {
	tmpFile, err := os.CreateTemp(dirPath, "tempname-*.tmp")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tmpFile.Close()
	return tmpFile.Name(), nil
}

// Relink removes file if it exists and
// creates a link at its place to the given path
func Relink(filePath, linkPath string) error {
	_ = os.Remove(filePath)
	return os.Link(linkPath, filePath)
}
