package saves

import (
	"errors"
	"os"
	"fmt"
	"regexp"
	"slices"
)

const (
	MAX_CHAPTER = 5
)

var (
	ErrOSNotSupportedHardLink = errors.New("Your OS is not supported due to hard links handling")
	ErrEmptySaveName          = errors.New("Empty save name is passed which is not allowed")
	ErrSaveNameIsTaken        = errors.New("Given save name is already taken")
	ErrChapterNotSupported    = errors.New("Given chapter is not yet supported")
	ErrSaveNotExist           = errors.New("At least on of the given saves doesn't exist")
	ErrShortSaveFile          = errors.New("Save file doesn't contain enough lines for parsing")
	ErrValueCannotBeSet       = errors.New("Can't set value during save parsing (Internal error)")
	ErrWrongLineType          = errors.New("Save file field contains value of an invalid type")
)

var (
	slotRegex = regexp.MustCompile(`filech(0|[^0\D]\d*)_(0|[^0\D]\d*)(_b)?`)
	saveRegex = regexp.MustCompile(`(0|[^0\D]\d*)_(a|b)_(.+)`)
)

func deleteEqual[S ~[]E, E comparable](s S, el E) S {
	return slices.DeleteFunc(s, func(x E) bool { return x == el })
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func tempFilePath(dirPath string) (string, error) {
	tmpFile, err := os.CreateTemp(dirPath, "tempname-*.tmp")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tmpFile.Close()
	return tmpFile.Name(), nil
}

func relink(filePath, linkPath string) error {
	_ = os.Remove(filePath)
	return os.Link(linkPath, filePath)
}
